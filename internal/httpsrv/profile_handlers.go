package httpsrv

import (
	"net/http"
	"strings"
)

// CVs are plain text; 60 KB ≈ a 10-page resume with headroom. The LLM prompt
// truncates far below this — the cap just keeps the column sane.
const maxCVBytes = 60 * 1024

func (s *Server) handleCVGet(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	var cv *string
	if err := s.Pool.QueryRow(r.Context(),
		`SELECT cv_text FROM users WHERE id = $1`, u.ID).Scan(&cv); err != nil {
		s.Logger.Error("cv read", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	out := ""
	if cv != nil {
		out = *cv
	}
	writeJSON(w, http.StatusOK, map[string]string{"cv_text": out})
}

// PUT /api/me/cv {cv_text} — empty string clears it.
func (s *Server) handleCVPut(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	var in struct {
		CVText string `json:"cv_text"`
	}
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	text := strings.TrimSpace(in.CVText)
	if len(text) > maxCVBytes {
		writeJSONError(w, http.StatusBadRequest, "CV is too long — trim it to the relevant parts (60 KB max)")
		return
	}
	var cvPtr *string
	if text != "" {
		cvPtr = &text
	}
	if _, err := s.Pool.Exec(r.Context(),
		`UPDATE users SET cv_text = $1 WHERE id = $2`, cvPtr, u.ID); err != nil {
		s.Logger.Error("cv save", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	s.logEvent(r.Context(), u.ID, "cv_save", map[string]any{"chars": len(text), "cleared": text == ""})
	writeJSON(w, http.StatusOK, map[string]any{"cv_chars": len(text)})
}

// POST /api/me/cv/extract {media_type, data} — PDF or image → plain text via
// Haiku. Returns the text for the user to review; nothing is saved until PUT.
func (s *Server) handleCVExtract(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	if s.LLM == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "AI features are not configured")
		return
	}
	var in parseImageReq
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	mt := strings.TrimSpace(in.MediaType)
	ok := mt == "application/pdf" || strings.HasPrefix(mt, "image/")
	if !ok || in.Data == "" {
		writeJSONError(w, http.StatusBadRequest, "upload a PDF or an image of your CV")
		return
	}
	text, err := s.LLM.ExtractCVText(r.Context(), mt, in.Data)
	if err != nil {
		s.logEvent(r.Context(), u.ID, "cv_extract", map[string]any{"outcome": "error"})
		writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	s.logEvent(r.Context(), u.ID, "cv_extract", map[string]any{"outcome": "ok", "chars": len(text)})
	writeJSON(w, http.StatusOK, map[string]string{"cv_text": text})
}

// PUT /api/me/reminders {enabled} — the in-app pair of the email unsubscribe link.
func (s *Server) handleRemindersPut(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	var in struct {
		Enabled bool `json:"enabled"`
	}
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := s.Pool.Exec(r.Context(),
		`UPDATE users SET email_reminders = $1 WHERE id = $2`, in.Enabled, u.ID); err != nil {
		s.Logger.Error("reminders toggle", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	s.logEvent(r.Context(), u.ID, "reminder_toggle", map[string]any{"enabled": in.Enabled})
	writeJSON(w, http.StatusOK, map[string]any{"email_reminders": in.Enabled})
}
