package httpsrv

import (
	"net/http"
	"strings"
	"time"
)

type authRequest struct {
	Email string `json:"email"`
}

func (s *Server) handleAuthRequest(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad json")
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || !strings.Contains(email, "@") {
		writeJSONError(w, http.StatusBadRequest, "invalid email")
		return
	}

	// Closed-beta gate. Respond with the same body either way so we don't leak
	// who's on the list.
	if !s.Cfg.EmailAllowed(email) {
		s.Logger.Info("auth.request blocked (not on allow-list)", "email", email)
		writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
		return
	}

	token, _, err := s.Auth.CreateMagicLink(r.Context(), email)
	if err != nil {
		s.Logger.Error("create magic link", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	link := s.Cfg.BaseURL + "/api/auth/verify?token=" + token
	if err := s.Mail.SendMagicLink(r.Context(), email, link); err != nil {
		s.Logger.Error("send magic link", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

func (s *Server) handleAuthVerify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSONError(w, http.StatusBadRequest, "missing token")
		return
	}
	session, u, err := s.Auth.ConsumeMagicLink(r.Context(), token, r.UserAgent(), clientIP(r))
	if err != nil {
		s.Logger.Info("magic link rejected", "err", err)
		// Redirect to login with an error flag so the UI can surface it.
		http.Redirect(w, r, "/?err=invalid_link", http.StatusFound)
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
	s.Logger.Info("user signed in", "user_id", u.ID, "email", u.Email)
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
		"id":    u.ID,
		"email": u.Email,
	})
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
