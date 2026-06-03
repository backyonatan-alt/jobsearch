package httpsrv

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type followUpDTO struct {
	ID            int64     `json:"id"`
	ApplicationID int64     `json:"application_id"`
	Note          string    `json:"note"`
	Channel       string    `json:"channel"`
	OccurredAt    time.Time `json:"occurred_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// GET /api/applications/{id}/follow-ups — newest first.
func (s *Server) handleFollowUpsList(w http.ResponseWriter, r *http.Request) {
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
		SELECT id, application_id, note, channel, occurred_at, created_at
		FROM follow_ups
		WHERE application_id = $1 AND user_id = $2
		ORDER BY occurred_at DESC, id DESC`, appID, u.ID)
	if err != nil {
		s.Logger.Error("follow-ups list", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()

	out := []followUpDTO{}
	for rows.Next() {
		var f followUpDTO
		if err := rows.Scan(&f.ID, &f.ApplicationID, &f.Note, &f.Channel,
			&f.OccurredAt, &f.CreatedAt); err != nil {
			s.Logger.Error("follow-up scan", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		out = append(out, f)
	}
	writeJSON(w, http.StatusOK, out)
}

type followUpCreateRequest struct {
	Note       string     `json:"note"`
	Channel    string     `json:"channel"`
	OccurredAt *time.Time `json:"occurred_at,omitempty"`
}

// POST /api/applications/{id}/follow-ups — log self-made outreach. Inserting a
// follow-up bumps the application's last_follow_up_at (computed in the list/get
// queries as MAX(occurred_at)), which is how "the clock resets".
func (s *Server) handleFollowUpCreate(w http.ResponseWriter, r *http.Request) {
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

	var in followUpCreateRequest
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	note := strings.TrimSpace(in.Note)
	channel := strings.TrimSpace(in.Channel)
	occurred := time.Now()
	if in.OccurredAt != nil && !in.OccurredAt.IsZero() {
		occurred = *in.OccurredAt
	}

	var f followUpDTO
	err = s.Pool.QueryRow(r.Context(), `
		INSERT INTO follow_ups (application_id, user_id, note, channel, occurred_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, application_id, note, channel, occurred_at, created_at`,
		appID, u.ID, note, channel, occurred,
	).Scan(&f.ID, &f.ApplicationID, &f.Note, &f.Channel, &f.OccurredAt, &f.CreatedAt)
	if err != nil {
		s.Logger.Error("follow-up insert", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusCreated, f)
}

// DELETE /api/applications/{id}/follow-ups/{fid}
func (s *Server) handleFollowUpDelete(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	appID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	fid, err := strconv.ParseInt(r.PathValue("fid"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad follow-up id")
		return
	}
	tag, err := s.Pool.Exec(r.Context(),
		`DELETE FROM follow_ups WHERE id = $1 AND application_id = $2 AND user_id = $3`,
		fid, appID, u.ID)
	if err != nil {
		s.Logger.Error("follow-up delete", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if tag.RowsAffected() == 0 {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
