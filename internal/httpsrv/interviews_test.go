package httpsrv

import "testing"

// Regression: AI-parsed events (free-text/screenshot via Haiku) carry
// source="ai". They were rejected at save with "source must be 'ics' or
// 'manual'", so a parsed interview could be previewed but never saved.
func TestValidInterviewSource(t *testing.T) {
	for _, s := range []string{"ics", "ai", "manual"} {
		if !validInterviewSource(s) {
			t.Errorf("expected %q to be a valid interview source", s)
		}
	}
	for _, s := range []string{"", "AI", "auto", "calendar", "llm"} {
		if validInterviewSource(s) {
			t.Errorf("expected %q to be rejected", s)
		}
	}
}
