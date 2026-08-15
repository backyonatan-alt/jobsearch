package httpsrv

import (
	"net/http"
	"strconv"

	"github.com/backyonatan-alt/jobsearch/internal/mail"
)

// GET /unsubscribe?u=<user_id>&t=<hmac> — one-click opt-out from reminder
// emails. Signed link, no session required (opened from a mail client).
func (s *Server) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.URL.Query().Get("u"), 10, 64)
	token := r.URL.Query().Get("t")
	if err != nil || !mail.ValidUnsubToken(s.Cfg.GoogleClientSecret, userID, token) {
		http.Error(w, "invalid link", http.StatusBadRequest)
		return
	}
	if _, err := s.Pool.Exec(r.Context(), `
		UPDATE users SET email_reminders = false WHERE id = $1`, userID); err != nil {
		s.Logger.Error("unsubscribe", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	s.logEvent(r.Context(), userID, "reminder_unsubscribe", nil)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!doctype html><meta name="viewport" content="width=device-width,initial-scale=1">
<div style="font-family:-apple-system,Segoe UI,Helvetica,Arial,sans-serif;max-width:480px;margin:80px auto;padding:0 24px;color:#1a1a2e">
<h2 style="font-size:19px">You're unsubscribed</h2>
<p style="color:#374151">Pursuit won't send you interview reminders anymore. Changed your mind? Reply to any Pursuit email and we'll switch them back on.</p>
<a href="/app" style="color:#1d4ed8">Back to Pursuit &rarr;</a></div>`))
}
