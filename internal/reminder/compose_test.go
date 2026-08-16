package reminder

import (
	"strings"
	"testing"
	"time"

	"github.com/backyonatan-alt/jobsearch/internal/mail"
)

func TestComposePreRoundTomorrow(t *testing.T) {
	tz := "Asia/Jerusalem"
	d := dueRound{
		InterviewID: 7, UserID: 3, ApplicationID: 5, Email: "x@y.z",
		Summary: "System design", Company: "Acme <Co>", Role: "Backend Engineer",
		StartsAt: time.Date(2026, 8, 20, 11, 0, 0, 0, time.UTC), TZ: &tz,
	}
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	subject, html, text := composePreRound("https://pursuit-playbook.com/", "sekret", d, now)

	if subject != "Your System design at Acme <Co> is tomorrow" {
		t.Fatalf("subject = %q", subject)
	}
	// 11:00 UTC in August Israel time (UTC+3) → Thursday 14:00.
	if !strings.Contains(html, "Thursday, Aug 20 · 14:00") {
		t.Errorf("html missing local time: %s", html)
	}
	wantLink := "https://pursuit-playbook.com/app/5/brief/7?src=reminder&utm_source=reminder&utm_medium=email"
	if !strings.Contains(html, wantLink) || !strings.Contains(text, wantLink) {
		t.Errorf("missing tagged brief link")
	}
	wantUnsub := "/unsubscribe?u=3&t=" + mail.UnsubToken("sekret", 3)
	if !strings.Contains(html, wantUnsub) || !strings.Contains(text, wantUnsub) {
		t.Errorf("missing unsubscribe link")
	}
	if !strings.Contains(html, "Acme &lt;Co&gt;") {
		t.Errorf("company not escaped in html")
	}
}

func TestComposePreRoundTodayFallbackTZ(t *testing.T) {
	d := dueRound{
		InterviewID: 1, UserID: 2, ApplicationID: 4, Email: "x@y.z",
		Summary: "", Company: "Wix", Role: "PM",
		StartsAt: time.Date(2026, 8, 19, 15, 0, 0, 0, time.UTC), TZ: nil,
	}
	now := time.Date(2026, 8, 19, 6, 0, 0, 0, time.UTC)
	subject, html, _ := composePreRound("https://pursuit-playbook.com", "sekret", d, now)
	if subject != "Your interview at Wix is today" {
		t.Fatalf("subject = %q", subject)
	}
	// Fallback zone is Asia/Jerusalem → 18:00.
	if !strings.Contains(html, "18:00") {
		t.Errorf("html missing fallback-tz time: %s", html)
	}
}

func TestUnsubTokenRoundTrip(t *testing.T) {
	tok := mail.UnsubToken("sekret", 42)
	if !mail.ValidUnsubToken("sekret", 42, tok) {
		t.Fatal("valid token rejected")
	}
	if mail.ValidUnsubToken("sekret", 43, tok) || mail.ValidUnsubToken("other", 42, tok) {
		t.Fatal("forged token accepted")
	}
}
