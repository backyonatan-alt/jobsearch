package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const resendEndpoint = "https://api.resend.com/emails"

// Client sends transactional email through the Resend HTTP API.
type Client struct {
	apiKey string
	from   string
	http   *http.Client
}

// New returns nil when apiKey is empty — callers treat a nil client as
// "email disabled" (same pattern as llm.New).
func New(apiKey, from string) *Client {
	if apiKey == "" {
		return nil
	}
	return &Client{
		apiKey: apiKey,
		from:   from,
		http:   &http.Client{Timeout: 15 * time.Second},
	}
}

type sendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html,omitempty"`
	Text    string   `json:"text,omitempty"`
}

func (c *Client) Send(ctx context.Context, to, subject, html, text string) error {
	body, err := json.Marshal(sendRequest{
		From: c.from, To: []string{to}, Subject: subject, HTML: html, Text: text,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("resend: %s: %s", resp.Status, string(b))
	}
	return nil
}
