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
	"regexp"
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
// arbitrary text. If the input is a bare URL we fetch it server-side first
// (LinkedIn is rejected since they block scraping). Errors include the raw
// model output on parse failure, which is useful for debugging prompt drift.
func (c *Client) ParseJob(ctx context.Context, text string) (*ParsedJob, error) {
	text = strings.TrimSpace(text)

	// Bare-URL path: fetch the page server-side so Claude has something to read.
	if u, ok := isBareURL(text); ok {
		if isLinkedInURL(u) {
			return nil, errors.New("LinkedIn blocks page fetching — copy the JD body text from the LinkedIn page and paste that instead")
		}
		body, err := c.fetchURL(ctx, u)
		if err != nil {
			return nil, fmt.Errorf("couldn't fetch %s (%v) — try pasting the JD text directly", u, err)
		}
		text = "Source URL: " + u + "\n\nPage content:\n" + body
	}

	resp, err := c.CreateMessage(ctx, ModelHaiku, parseSystemPrompt, []Message{
		{Role: "user", Content: text},
	}, 600)
	if err != nil {
		return nil, err
	}

	raw, err := extractFirstJSONObject(resp)
	if err != nil {
		return nil, fmt.Errorf("model didn't return JSON (raw: %.200s)", resp)
	}

	// Explicit-error reply: surface a friendlier message.
	var errProbe struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(raw, &errProbe) == nil && errProbe.Error != "" {
		if strings.EqualFold(errProbe.Error, "not a job listing") {
			return nil, errors.New("that doesn't look like a job listing — try pasting the JD body text instead")
		}
		return nil, errors.New(errProbe.Error)
	}

	var job ParsedJob
	if err := json.Unmarshal(raw, &job); err != nil {
		return nil, fmt.Errorf("parse model JSON: %w", err)
	}
	if job.Company == "" && job.Role == "" {
		return nil, errors.New("could not extract company or role from the input")
	}
	return &job, nil
}

var (
	bareURLRe = regexp.MustCompile(`^https?://\S+$`)
	scriptRe  = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	styleRe   = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	tagRe     = regexp.MustCompile(`<[^>]+>`)
	wsRe      = regexp.MustCompile(`\s+`)
)

func isBareURL(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if !bareURLRe.MatchString(s) || len(s) > 500 {
		return "", false
	}
	return s, true
}

func isLinkedInURL(u string) bool {
	lower := strings.ToLower(u)
	return strings.Contains(lower, "linkedin.com/")
}

// fetchURL retrieves a public URL and returns visible text content (HTML tags,
// scripts, and styles stripped). Capped to ~30k chars so we don't blow the
// token budget on huge pages.
func (c *Client) fetchURL(ctx context.Context, url string) (string, error) {
	cctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(cctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Pursuit/0.1)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 400_000))
	if err != nil {
		return "", err
	}
	text := stripHTML(string(raw))
	if len(text) > 30_000 {
		text = text[:30_000]
	}
	return text, nil
}

func stripHTML(s string) string {
	s = scriptRe.ReplaceAllString(s, " ")
	s = styleRe.ReplaceAllString(s, " ")
	s = tagRe.ReplaceAllString(s, " ")
	s = wsRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// extractFirstJSONObject finds the first top-level JSON object in s and returns
// it as bytes, ignoring any prose, markdown fences, or trailing text the model
// may have wrapped it in.
func extractFirstJSONObject(s string) ([]byte, error) {
	start := strings.Index(s, "{")
	if start == -1 {
		return nil, errors.New("no JSON object found")
	}
	dec := json.NewDecoder(strings.NewReader(s[start:]))
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}
