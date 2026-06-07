package httpsrv

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/backyonatan-alt/jobsearch/internal/llm"
)

type parseImage struct {
	MediaType string `json:"media_type"`
	Data      string `json:"data"` // base64, no data: prefix
}

type parseRequest struct {
	Text  string      `json:"text"`
	Image *parseImage `json:"image,omitempty"`
}

// Anthropic accepts these four image types. Cap raw bytes at 5 MB per their
// docs; 5MB raw is ~6.7MB base64 — we accept up to 8MB encoded to leave slack.
var allowedImageMediaTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
	"image/webp": true,
}

const maxImageBase64Bytes = 8 * 1024 * 1024 // ~6 MB raw

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

	var img *llm.ParseImage
	if req.Image != nil && req.Image.Data != "" {
		if !allowedImageMediaTypes[req.Image.MediaType] {
			writeJSONError(w, http.StatusBadRequest, "image media_type must be image/png, image/jpeg, image/gif, or image/webp")
			return
		}
		if len(req.Image.Data) > maxImageBase64Bytes {
			writeJSONError(w, http.StatusBadRequest, "image too large (8 MB base64 max — try a smaller screenshot)")
			return
		}
		// Validate the base64 is decodable so we don't pay Anthropic to reject it.
		if _, err := base64.StdEncoding.DecodeString(req.Image.Data); err != nil {
			writeJSONError(w, http.StatusBadRequest, "image data is not valid base64")
			return
		}
		img = &llm.ParseImage{MediaType: req.Image.MediaType, Data: req.Image.Data}
	}

	if img == nil {
		if len(text) < 5 {
			writeJSONError(w, http.StatusBadRequest, "paste a job listing or URL, or attach a screenshot")
			return
		}
		if len(text) > 50_000 {
			writeJSONError(w, http.StatusBadRequest, "text too long (50k chars max)")
			return
		}
	} else if len(text) > 50_000 {
		writeJSONError(w, http.StatusBadRequest, "text too long (50k chars max)")
		return
	}

	// The screenshot path and the paste path are different product moments,
	// so they get distinct event names. Emitted server-side with true latency
	// so a client crash on error can't drop the signal.
	u, _ := userFromCtx(r.Context())
	eventName := "paste_parse"
	if img != nil {
		eventName = "screenshot_parse"
	}
	start := time.Now()

	job, err := s.LLM.ParseJob(r.Context(), text, img)
	if err != nil {
		s.Logger.Info("parse failed", "err", err, "has_image", img != nil)
		s.logEvent(r.Context(), u.ID, eventName, map[string]any{
			"outcome": "error", "error_reason": "parse_failed",
			"duration_ms": time.Since(start).Milliseconds(),
		})
		writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	s.logEvent(r.Context(), u.ID, eventName, map[string]any{
		"outcome": "success", "duration_ms": time.Since(start).Milliseconds(),
	})
	writeJSON(w, http.StatusOK, job)
}
