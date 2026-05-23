// Package llm wraps the Anthropic Messages API. We call the REST endpoint
// directly to avoid pulling in an SDK for what amounts to two endpoints
// (parse and dossier).
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	apiBase    = "https://api.anthropic.com"
	apiVersion = "2023-06-01"

	ModelHaiku  = "claude-haiku-4-5-20251001"
	ModelSonnet = "claude-sonnet-4-6"
)

// Client is safe for concurrent use.
type Client struct {
	apiKey string
	http   *http.Client
}

// New returns a Client, or nil if apiKey is empty (so callers can keep the
// LLM-dependent endpoints behind a "not configured" 503 in dev).
func New(apiKey string) *Client {
	if strings.TrimSpace(apiKey) == "" {
		return nil
	}
	return &Client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 60 * time.Second},
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messagesRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
	Messages  []Message `json:"messages"`
}

type messagesResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// CreateMessage hits POST /v1/messages and returns the concatenated text
// content across all text blocks in the assistant response.
func (c *Client) CreateMessage(ctx context.Context, model, system string, messages []Message, maxTokens int) (string, error) {
	body, err := json.Marshal(messagesRequest{
		Model:     model,
		MaxTokens: maxTokens,
		System:    system,
		Messages:  messages,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", apiBase+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", apiVersion)
	req.Header.Set("content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic %d: %s", resp.StatusCode, string(raw))
	}
	var out messagesResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("parse anthropic response: %w", err)
	}
	if len(out.Content) == 0 {
		return "", errors.New("anthropic returned empty content")
	}
	var text strings.Builder
	for _, b := range out.Content {
		if b.Type == "text" {
			text.WriteString(b.Text)
		}
	}
	return text.String(), nil
}

// ParsedJob is the structured result we expect Claude to extract from arbitrary
// pasted text. All fields are best-effort; omitted fields mean "not found".
type ParsedJob struct {
	Company    string `json:"company,omitempty"`
	Role       string `json:"role,omitempty"`
	Location   string `json:"location,omitempty"`
	Seniority  string `json:"seniority,omitempty"`
	JDURL      string `json:"jd_url,omitempty"`
	Source     string `json:"source,omitempty"`
	SalaryNote string `json:"salary_note,omitempty"`
}

const parseSystemPrompt = `You are a precise job-listing parser. The user will paste arbitrary text — a LinkedIn URL, a copied job description, an email forward, a recruiter message, or a description in their own words — and you must extract structured fields.

Return ONLY a JSON object with these keys (omit any you cannot determine with high confidence; do not invent values):
- company:    the hiring company name (e.g. "Anthropic")
- role:       the role title as listed (e.g. "Senior Software Engineer")
- location:   location or remote status (e.g. "San Francisco", "Remote", "Remote (US)")
- seniority:  level if clear (e.g. "Senior", "Staff", "Principal", "Director")
- jd_url:     the canonical job-description URL if present in the text
- source:     "LinkedIn" if from a LinkedIn URL or page, else infer (e.g. "Greenhouse", "Lever", "Company site", "Referral", "Email") or omit
- salary_note: compensation info if mentioned (e.g. "$220k-$280k base", "$160-200k OTE")

Output the JSON object only, no prose, no markdown fences. If the input is clearly not a job listing, return {"error": "not a job listing"}.`

// ParseJob asks Claude Haiku to extract structured job-listing fields from
// arbitrary text. Errors include the raw model output on parse failure, which
// is useful for debugging prompt drift.
func (c *Client) ParseJob(ctx context.Context, text string) (*ParsedJob, error) {
	resp, err := c.CreateMessage(ctx, ModelHaiku, parseSystemPrompt, []Message{
		{Role: "user", Content: text},
	}, 600)
	if err != nil {
		return nil, err
	}
	cleaned := stripFences(strings.TrimSpace(resp))

	// First check for an explicit error reply.
	var errProbe struct {
		Error string `json:"error"`
	}
	if json.Unmarshal([]byte(cleaned), &errProbe) == nil && errProbe.Error != "" {
		return nil, errors.New(errProbe.Error)
	}

	var job ParsedJob
	if err := json.Unmarshal([]byte(cleaned), &job); err != nil {
		return nil, fmt.Errorf("parse model JSON: %w (raw: %s)", err, cleaned)
	}
	if job.Company == "" && job.Role == "" {
		return nil, errors.New("could not extract company or role from the input")
	}
	return &job, nil
}

// stripFences removes markdown code fences if the model wrapped its JSON in them
// despite being told not to.
func stripFences(s string) string {
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
