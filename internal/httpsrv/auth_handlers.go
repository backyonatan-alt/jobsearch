package httpsrv

import (
	"net/http"
	"strings"
	"time"

	"github.com/backyonatan-alt/jobsearch/internal/auth"
)

const oauthStateCookie = "pursuit_oauth_state"

func (s *Server) handleGoogleStart(w http.ResponseWriter, r *http.Request) {
	state, err := auth.RandomToken(24)
	if err != nil {
		s.Logger.Error("oauth state", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    state,
		Path:     "/auth/google/",
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   strings.HasPrefix(s.Cfg.BaseURL, "https://"),
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, s.Google.AuthCodeURL(state), http.StatusFound)
}

func (s *Server) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if oauthErr := q.Get("error"); oauthErr != "" {
		s.Logger.Info("google oauth error", "err", oauthErr)
		http.Redirect(w, r, "/?err=oauth_denied", http.StatusFound)
		return
	}

	stateCookie, err := r.Cookie(oauthStateCookie)
	if err != nil || stateCookie.Value == "" || stateCookie.Value != q.Get("state") {
		s.Logger.Info("oauth state mismatch")
		http.Redirect(w, r, "/?err=oauth_state", http.StatusFound)
		return
	}
	// State cookie consumed — clear it.
	http.SetCookie(w, &http.Cookie{
		Name:    oauthStateCookie,
		Value:   "",
		Path:    "/auth/google/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	code := q.Get("code")
	if code == "" {
		http.Redirect(w, r, "/?err=oauth_no_code", http.StatusFound)
		return
	}

	email, err := s.Google.ExchangeAndVerify(r.Context(), code)
	if err != nil {
		s.Logger.Info("google exchange failed", "err", err)
		http.Redirect(w, r, "/?err=oauth_failed", http.StatusFound)
		return
	}

	invited, err := s.Auth.EmailInvited(r.Context(), email, s.Cfg.AllowedEmails)
	if err != nil {
		s.Logger.Error("invite check", "err", err)
		http.Redirect(w, r, "/?err=internal", http.StatusFound)
		return
	}
	if !invited {
		s.Logger.Info("oauth blocked (not on allow-list)", "email", email)
		http.Redirect(w, r, "/?err=not_invited", http.StatusFound)
		return
	}

	session, u, err := s.Auth.UpsertUserAndSession(r.Context(), email, r.UserAgent(), clientIP(r), s.Cfg.IsAdminEmail(email))
	if err != nil {
		s.Logger.Error("upsert user", "err", err)
		http.Redirect(w, r, "/?err=internal", http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.Cfg.SessionCookieName,
		Value:    session,
		Path:     "/",
		Expires:  time.Now().Add(s.Cfg.SessionTTL),
		HttpOnly: true,
		Secure:   strings.HasPrefix(s.Cfg.BaseURL, "https://"),
		SameSite: http.SameSiteLaxMode,
	})
	s.Logger.Info("user signed in via google", "user_id", u.ID, "email", u.Email)
	http.Redirect(w, r, "/app", http.StatusFound)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(s.Cfg.SessionCookieName); err == nil {
		_ = s.Auth.DestroySession(r.Context(), c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     s.Cfg.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	writeJSON(w, http.StatusOK, map[string]any{
		"id":           u.ID,
		"email":        u.Email,
		"is_admin":     u.IsAdmin,
		"onboarded_at": u.OnboardedAt,
	})
}

func (s *Server) handleMarkOnboarded(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	if err := s.Auth.MarkOnboarded(r.Context(), u.ID); err != nil {
		s.Logger.Error("mark onboarded", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if i := strings.IndexByte(fwd, ','); i >= 0 {
			return strings.TrimSpace(fwd[:i])
		}
		return strings.TrimSpace(fwd)
	}
	host := r.RemoteAddr
	if i := strings.LastIndexByte(host, ':'); i > 0 {
		host = host[:i]
	}
	return host
}
