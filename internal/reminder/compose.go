package reminder

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/backyonatan-alt/jobsearch/internal/mail"
)

// Closed beta skews Israeli; used only when the browser never told us a zone.
const fallbackTZ = "Asia/Jerusalem"

func composePreRound(baseURL, unsubSecret string, d dueRound, now time.Time) (subject, htmlBody, text string) {
	loc, err := time.LoadLocation(fallbackTZ)
	if d.TZ != nil {
		if l, e := time.LoadLocation(*d.TZ); e == nil {
			loc, err = l, nil
		}
	}
	if err != nil {
		loc = time.UTC
	}
	starts := d.StartsAt.In(loc)

	day := "coming up"
	nowDay := now.In(loc).Format("2006-01-02")
	switch starts.Format("2006-01-02") {
	case nowDay:
		day = "today"
	case now.In(loc).Add(24 * time.Hour).Format("2006-01-02"):
		day = "tomorrow"
	}

	when := starts.Format("Monday, Jan 2")
	if !d.AllDay {
		when += " · " + starts.Format("15:04")
	}

	round := strings.TrimSpace(d.Summary)
	if round == "" {
		round = "interview"
	}
	subject = fmt.Sprintf("Your %s at %s is %s", round, d.Company, day)

	base := strings.TrimRight(baseURL, "/")
	briefURL := fmt.Sprintf("%s/app/%d/brief/%d?src=reminder&utm_source=reminder&utm_medium=email",
		base, d.ApplicationID, d.InterviewID)
	unsubURL := fmt.Sprintf("%s/unsubscribe?u=%d&t=%s", base, d.UserID, mail.UnsubToken(unsubSecret, d.UserID))

	esc := html.EscapeString
	htmlBody = fmt.Sprintf(`<div style="font-family:-apple-system,Segoe UI,Helvetica,Arial,sans-serif;max-width:520px;margin:0 auto;padding:24px;color:#1a1a2e">
  <p style="font-size:13px;color:#6b7280;margin:0 0 16px">Pursuit</p>
  <h2 style="font-size:19px;margin:0 0 8px">Your %s at %s is %s</h2>
  <p style="margin:0 0 4px;color:#374151">%s</p>
  <p style="margin:0 0 20px;color:#374151">%s — %s</p>
  <a href="%s" style="display:inline-block;background:#1d4ed8;color:#fff;text-decoration:none;padding:10px 18px;border-radius:8px;font-size:15px">Open your playbook &rarr;</a>
  <p style="margin:28px 0 0;font-size:12px;color:#9ca3af">You're getting this because a round is scheduled in Pursuit.
  <a href="%s" style="color:#9ca3af">Stop interview reminders</a></p>
</div>`,
		esc(round), esc(d.Company), day, esc(when), esc(d.Role), esc(d.Company), briefURL, unsubURL)

	text = fmt.Sprintf("Your %s at %s is %s\n%s\n%s — %s\n\nOpen your playbook: %s\n\nStop interview reminders: %s\n",
		round, d.Company, day, when, d.Role, d.Company, briefURL, unsubURL)
	return subject, htmlBody, text
}
