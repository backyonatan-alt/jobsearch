package httpsrv

import (
	"net/http"
	"strings"
)

type parseRequest struct {
	Text string `json:"text"`
}

func (s *Server) handleApplicationParse(w http.ResponseWriter, r *http.Request) {
	if s.LLM == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "AI parsing is not configured (ANTHROPIC_API_KEY missing)")
		return
	}
	var req parseRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	text := strings.TrimSpace(req.Text)
	if len(text) < 5 {
		writeJSONError(w, http.StatusBadRequest, "text too short")
		return
	}
	if len(text) > 50_000 {
		writeJSONError(w, http.StatusBadRequest, "text too long (50k chars max)")
		return
	}

	job, err := s.LLM.ParseJob(r.Context(), text)
	if err != nil {
		s.Logger.Info("parse failed", "err", err)
		writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, job)
}
