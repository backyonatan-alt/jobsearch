package llm

import (
	"encoding/json"
	"strings"
	"testing"
)

// Anthropic's Messages API accepts content as EITHER a string or an array of
// content blocks. We rely on Go's `any` field serializing both shapes
// correctly. This test pins the wire format so a future refactor can't
// silently break vision (`content: "..."`-only would drop the image).
func TestMessageContentSerialization_String(t *testing.T) {
	m := Message{Role: "user", Content: "hello"}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := `{"role":"user","content":"hello"}`
	if got != want {
		t.Errorf("text message:\n got: %s\nwant: %s", got, want)
	}
}

func TestMessageContentSerialization_ImageBlocks(t *testing.T) {
	m := Message{
		Role: "user",
		Content: []ContentBlock{
			ImageBlock("image/png", "AAAA"),
			TextBlock("Extract company and role."),
		},
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := `{"role":"user","content":[` +
		`{"type":"image","source":{"type":"base64","media_type":"image/png","data":"AAAA"}},` +
		`{"type":"text","text":"Extract company and role."}` +
		`]}`
	if got != want {
		t.Errorf("image message:\n got: %s\nwant: %s", got, want)
	}
}

// TextBlock must NOT emit a null `source` field — Anthropic rejects text
// blocks that include source. omitempty on the pointer is what saves us.
func TestTextBlockOmitsSource(t *testing.T) {
	b, err := json.Marshal(TextBlock("hi"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "source") {
		t.Errorf("TextBlock leaked source field: %s", b)
	}
}

// The round's format (e.g. "Home Assignment Presentation") must reach the
// model — a round with a stated format that generates generic prep is the bug
// this pins against. Also pins that an empty round context adds no line.
func TestInterviewerPromptCarriesRoundContext(t *testing.T) {
	got := buildInterviewerPromptUser(
		"HiBob", "Senior Product Manager - Core HR", "Adir Nashawi",
		"Tel Aviv", "https://www.hibob.com",
		"Home Assignment Presentation — walk through the take-home and defend it",
		"round 1 went strong",
	)
	for _, want := range []string{
		"Company: HiBob",
		"Interviewer: Adir Nashawi",
		"This round: Home Assignment Presentation — walk through the take-home and defend it",
		"round 1 went strong",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("prompt missing %q:\n%s", want, got)
		}
	}

	plain := buildInterviewerPromptUser("HiBob", "PM", "Adir", "", "", "", "")
	if strings.Contains(plain, "This round:") {
		t.Errorf("empty round context should add no line:\n%s", plain)
	}
}
