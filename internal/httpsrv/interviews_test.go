package httpsrv

import (
	"encoding/json"
	"testing"
	"time"
)

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

// Golden round-trip: the parse endpoint returns parsedEventDTO; the frontend
// posts those objects straight back to the save endpoint, which reads them as
// interviewCreateRequest. If the two shapes drift (a renamed json tag, a
// dropped time field), a parsed interview previews but never saves — exactly
// the 0/22 hole we found. Assert the full object survives the round-trip and
// would pass the save handler's validation.
func TestParsedEventSavesRoundTrip(t *testing.T) {
	start := time.Date(2026, 6, 25, 14, 30, 0, 0, time.UTC)
	end := start.Add(45 * time.Minute)
	parsed := parsedEventDTO{
		Source:      "ai",
		Summary:     "Hiring manager interview",
		Location:    "Google Meet",
		Description: "45 min with Dana",
		StartsAt:    start,
		EndsAt:      &end,
		AllDay:      false,
		Attendees:   []person{{Name: "Dana Levi"}, {Email: "dana@acme.com"}},
	}

	wire, err := json.Marshal(parsed)
	if err != nil {
		t.Fatalf("marshal parsed event: %v", err)
	}
	var in interviewCreateRequest
	if err := json.Unmarshal(wire, &in); err != nil {
		t.Fatalf("unmarshal into create request: %v", err)
	}

	// The two failure modes that produce "previewed but never saved".
	if !validInterviewSource(in.Source) {
		t.Errorf("round-tripped source %q would be rejected at save", in.Source)
	}
	if in.StartsAt.IsZero() {
		t.Error("starts_at did not survive the round-trip — save would 400 'starts_at is required'")
	}

	// Everything else should come through intact.
	if !in.StartsAt.Equal(start) {
		t.Errorf("starts_at = %v, want %v", in.StartsAt, start)
	}
	if in.EndsAt == nil || !in.EndsAt.Equal(end) {
		t.Errorf("ends_at = %v, want %v", in.EndsAt, end)
	}
	if in.Summary != parsed.Summary || in.Location != parsed.Location || in.Description != parsed.Description {
		t.Errorf("text fields drifted: %+v", in)
	}
	if len(in.Attendees) != 2 || in.Attendees[0].Name != "Dana Levi" || in.Attendees[1].Email != "dana@acme.com" {
		t.Errorf("attendees drifted: %+v", in.Attendees)
	}
}
