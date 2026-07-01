package httpsrv

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type debriefDTO struct {
	ID           int64     `json:"id"`
	InterviewID  int64     `json:"interview_id"`
	Feel         string    `json:"feel"`
	PrepAccuracy string    `json:"prep_accuracy"`
	Topics       string    `json:"topics"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type debriefInput struct {
	Feel         string `json:"feel"`
	PrepAccuracy string `json:"prep_accuracy"`
	Topics       string `json:"topics"`
	Notes        string `json:"notes"`
}

var validFeel = map[string]bool{"strong": true, "mixed": true, "rough": true}
var validPrepAccuracy = map[string]bool{"spot_on": true, "partly": true, "off": true}

// GET /api/applications/{id}/debriefs — all debriefs for the application, so the
// detail page knows which rounds are debriefed (for the capture card + the
// "informed by your last round" chip) in one fetch.
func (s *Server) handleDebriefsList(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	appID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	if _, err := s.fetchApplication(r.Context(), u.ID, appID); err != nil {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}
	rows, err := s.Pool.Query(r.Context(), `
		SELECT id, interview_id, feel, prep_accuracy, topics, notes, created_at, updated_at
		FROM debriefs WHERE application_id = $1 AND user_id = $2
		ORDER BY created_at ASC`, appID, u.ID)
	if err != nil {
		s.Logger.Error("debriefs list", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()
	out := []debriefDTO{}
	for rows.Next() {
		var d debriefDTO
		if err := rows.Scan(&d.ID, &d.InterviewID, &d.Feel, &d.PrepAccuracy, &d.Topics, &d.Notes, &d.CreatedAt, &d.UpdatedAt); err != nil {
			s.Logger.Error("debriefs scan", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		out = append(out, d)
	}
	writeJSON(w, http.StatusOK, out)
}

// POST /api/applications/{id}/interviews/{iid}/debrief — upsert one round's debrief.
func (s *Server) handleDebriefSave(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	appID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	iid, err := strconv.ParseInt(r.PathValue("iid"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad interview id")
		return
	}
	if _, err := s.fetchApplication(r.Context(), u.ID, appID); err != nil {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}
	iv, err := s.interviewByID(r.Context(), appID, iid)
	if err != nil {
		s.Logger.Error("debrief interview lookup", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if iv == nil {
		writeJSONError(w, http.StatusNotFound, "interview not found")
		return
	}

	var in debriefInput
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !validFeel[in.Feel] {
		writeJSONError(w, http.StatusBadRequest, "invalid feel")
		return
	}
	if !validPrepAccuracy[in.PrepAccuracy] {
		writeJSONError(w, http.StatusBadRequest, "invalid prep_accuracy")
		return
	}

	var d debriefDTO
	err = s.Pool.QueryRow(r.Context(), `
		INSERT INTO debriefs (interview_id, application_id, user_id, feel, prep_accuracy, topics, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (interview_id) DO UPDATE SET
		    feel = EXCLUDED.feel, prep_accuracy = EXCLUDED.prep_accuracy,
		    topics = EXCLUDED.topics, notes = EXCLUDED.notes, updated_at = now()
		RETURNING id, interview_id, feel, prep_accuracy, topics, notes, created_at, updated_at`,
		iid, appID, u.ID, in.Feel, in.PrepAccuracy, strings.TrimSpace(in.Topics), strings.TrimSpace(in.Notes),
	).Scan(&d.ID, &d.InterviewID, &d.Feel, &d.PrepAccuracy, &d.Topics, &d.Notes, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		s.Logger.Error("debrief upsert", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusOK, d)
	// Trust signal — the raw material for "was the prep right?".
	s.logEvent(r.Context(), u.ID, "debrief_save", map[string]any{"feel": in.Feel, "prep_accuracy": in.PrepAccuracy})
}

// priorDebriefsContext assembles a short natural-language summary of debriefs
// from rounds that happened BEFORE `before`, so the next round's interviewer
// brief can build on what already came up. Empty string when there's nothing.
func (s *Server) priorDebriefsContext(ctx context.Context, appID int64, before time.Time) string {
	rows, err := s.Pool.Query(ctx, `
		SELECT i.summary, i.starts_at, d.feel, d.prep_accuracy, d.topics, d.notes
		FROM debriefs d JOIN interviews i ON i.id = d.interview_id
		WHERE d.application_id = $1 AND i.starts_at < $2
		ORDER BY i.starts_at ASC`, appID, before)
	if err != nil {
		return ""
	}
	defer rows.Close()
	var b strings.Builder
	for rows.Next() {
		var summary, feel, acc, topics, notes string
		var startsAt time.Time
		if rows.Scan(&summary, &startsAt, &feel, &acc, &topics, &notes) != nil {
			continue
		}
		label := strings.TrimSpace(summary)
		if label == "" {
			label = "an earlier round"
		}
		fmt.Fprintf(&b, "- %s (%s): it went %s; the candidate rated our prep %s.",
			label, startsAt.Format("Jan 2"), feelWord(feel), accWord(acc))
		if topics != "" {
			fmt.Fprintf(&b, " What actually came up: %s.", topics)
		}
		if notes != "" {
			fmt.Fprintf(&b, " Their note: %s.", notes)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func feelWord(f string) string {
	switch f {
	case "strong":
		return "well"
	case "mixed":
		return "so-so"
	case "rough":
		return "rough"
	}
	return "unclear"
}

func accWord(a string) string {
	switch a {
	case "spot_on":
		return "spot-on"
	case "partly":
		return "partly right"
	case "off":
		return "off"
	}
	return "unrated"
}
