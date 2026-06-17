package httpsrv

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/backyonatan-alt/jobsearch/internal/ics"
	"github.com/backyonatan-alt/jobsearch/internal/llm"
	"github.com/jackc/pgx/v5"
)

// Cap raw .ics payloads at 256 KB. Real-world invites are under 10 KB; this
// leaves headroom for huge DESCRIPTION blocks (calendar attachments inline as
// base64 in DESCRIPTION sometimes) without letting anyone post megabytes.
const maxICSBytes = 256 * 1024

type interviewDTO struct {
	ID            int64           `json:"id"`
	ApplicationID int64           `json:"application_id"`
	Source        string          `json:"source"`
	UID           *string         `json:"uid,omitempty"`
	Summary       string          `json:"summary"`
	Location      *string         `json:"location,omitempty"`
	Description   *string         `json:"description,omitempty"`
	StartsAt      time.Time       `json:"starts_at"`
	EndsAt        *time.Time      `json:"ends_at,omitempty"`
	AllDay        bool            `json:"all_day"`
	Organizer     json.RawMessage `json:"organizer,omitempty"`
	Attendees     json.RawMessage `json:"attendees"`
	CreatedAt     time.Time       `json:"created_at"`
}

type parsedEventDTO struct {
	Source      string     `json:"source"`
	UID         string     `json:"uid"`
	Summary     string     `json:"summary"`
	Location    string     `json:"location"`
	Description string     `json:"description"`
	StartsAt    time.Time  `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
	AllDay      bool       `json:"all_day"`
	Organizer   *person    `json:"organizer,omitempty"`
	Attendees   []person   `json:"attendees"`
}

type person struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type icsParseRequest struct {
	ICS   string         `json:"ics,omitempty"`
	Text  string         `json:"text,omitempty"`
	Image *parseImageReq `json:"image,omitempty"`
	TZ    string         `json:"tz,omitempty"` // IANA name from the browser, e.g. "Asia/Jerusalem"
}

type parseImageReq struct {
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type icsParseResponse struct {
	Events []parsedEventDTO `json:"events"`
}

// POST /api/applications/{id}/interviews/parse — three ways to parse:
//  1. {ics: "BEGIN:VCALENDAR..."}   — strict ICS parsing, fastest, no LLM.
//  2. {text: "You're invited..."}   — free-form email/notes through Haiku.
//  3. {image: {media_type, data}}   — screenshot through Haiku Vision.
//
// Returns the same {events: [...]} shape regardless so the frontend renders
// the preview through one code path.
func (s *Server) handleInterviewsParse(w http.ResponseWriter, r *http.Request) {
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
	start := time.Now()

	var req icsParseRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	icsBody := strings.TrimSpace(req.ICS)
	text := strings.TrimSpace(req.Text)
	if icsBody == "" && text == "" && req.Image == nil {
		writeJSONError(w, http.StatusBadRequest, "paste a calendar file, a screenshot, or the event text")
		return
	}

	// 1. ICS path — strict parsing, no LLM.
	if icsBody != "" {
		if len(icsBody) > maxICSBytes {
			writeJSONError(w, http.StatusBadRequest, "calendar payload too large (256 KB max)")
			return
		}
		events, err := ics.Parse(strings.NewReader(icsBody))
		if err != nil {
			s.logEvent(r.Context(), u.ID, "interview_parse", map[string]any{
				"source": "ics", "outcome": "error", "error_reason": "parse_failed",
				"duration_ms": time.Since(start).Milliseconds(),
			})
			writeJSONError(w, http.StatusUnprocessableEntity, "could not parse calendar: "+err.Error())
			return
		}
		if len(events) == 0 {
			s.logEvent(r.Context(), u.ID, "interview_parse", map[string]any{
				"source": "ics", "outcome": "error", "error_reason": "no_vevent",
				"duration_ms": time.Since(start).Milliseconds(),
			})
			writeJSONError(w, http.StatusUnprocessableEntity, "no VEVENT found — is this an .ics calendar invite?")
			return
		}
		out := make([]parsedEventDTO, 0, len(events))
		for _, e := range events {
			out = append(out, toParsedDTO(e))
		}
		s.logEvent(r.Context(), u.ID, "interview_parse", map[string]any{
			"source": "ics", "outcome": "success", "count": len(out),
			"duration_ms": time.Since(start).Milliseconds(),
		})
		writeJSON(w, http.StatusOK, icsParseResponse{Events: out})
		return
	}

	// 2/3. Image + text path — needs the LLM client.
	if s.LLM == nil {
		writeJSONError(w, http.StatusServiceUnavailable,
			"AI event parsing is not configured (ANTHROPIC_API_KEY missing) — paste an .ics file instead")
		return
	}
	var img *llm.ParseImage
	if req.Image != nil {
		if req.Image.Data == "" || req.Image.MediaType == "" {
			writeJSONError(w, http.StatusBadRequest, "image needs both media_type and data")
			return
		}
		img = &llm.ParseImage{MediaType: req.Image.MediaType, Data: req.Image.Data}
	}
	s.Logger.Info("interview parse start",
		"user_id", u.ID, "app_id", appID,
		"has_image", img != nil, "text_chars", len(text))

	source := "text"
	if img != nil {
		source = "image"
	}
	ev, err := s.LLM.ParseEvent(r.Context(), text, img, strings.TrimSpace(req.TZ))
	if err != nil {
		s.Logger.Info("interview parse failed", "err", err)
		s.logEvent(r.Context(), u.ID, "interview_parse", map[string]any{
			"source": source, "outcome": "error", "error_reason": "parse_failed",
			"duration_ms": time.Since(start).Milliseconds(),
		})
		writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	dto, err := llmEventToDTO(ev)
	if err != nil {
		s.logEvent(r.Context(), u.ID, "interview_parse", map[string]any{
			"source": source, "outcome": "error", "error_reason": "bad_event",
			"duration_ms": time.Since(start).Milliseconds(),
		})
		writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	s.logEvent(r.Context(), u.ID, "interview_parse", map[string]any{
		"source": source, "outcome": "success", "count": 1,
		"duration_ms": time.Since(start).Milliseconds(),
	})
	writeJSON(w, http.StatusOK, icsParseResponse{Events: []parsedEventDTO{dto}})
}

// llmEventToDTO converts the Haiku-parsed event (ISO strings, free-form
// attendee names) into the parsedEventDTO shape the rest of the pipeline
// expects (parsed time.Time, person structs).
func llmEventToDTO(ev *llm.ParsedEvent) (parsedEventDTO, error) {
	d := parsedEventDTO{
		Source:      "ai",
		Summary:     ev.Summary,
		Location:    ev.Location,
		Description: ev.Description,
		AllDay:      ev.AllDay,
		Attendees:   []person{},
	}
	if ev.StartsAt != "" {
		t, err := parseFlexibleTime(ev.StartsAt)
		if err != nil {
			return d, errors.New("model returned a start time we couldn't parse: " + ev.StartsAt)
		}
		d.StartsAt = t
	} else {
		return d, errors.New("model didn't return a start time")
	}
	if ev.EndsAt != "" {
		t, err := parseFlexibleTime(ev.EndsAt)
		if err == nil {
			d.EndsAt = &t
		}
	}
	for _, a := range ev.Attendees {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if strings.Contains(a, "@") {
			d.Attendees = append(d.Attendees, person{Email: a})
		} else {
			d.Attendees = append(d.Attendees, person{Name: a})
		}
	}
	return d, nil
}

// parseFlexibleTime accepts the common ISO 8601 shapes the model is likely to
// return — full RFC3339, "YYYY-MM-DDTHH:MM:SS", date-only.
func parseFlexibleTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	for _, layout := range []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unrecognized time format")
}

func toParsedDTO(e ics.Event) parsedEventDTO {
	d := parsedEventDTO{
		Source:      "ics",
		UID:         e.UID,
		Summary:     e.Summary,
		Location:    e.Location,
		Description: e.Description,
		StartsAt:    e.StartsAt,
		AllDay:      e.AllDay,
		Attendees:   []person{},
	}
	if !e.EndsAt.IsZero() {
		t := e.EndsAt
		d.EndsAt = &t
	}
	if e.Organizer.Email != "" || e.Organizer.Name != "" {
		d.Organizer = &person{Name: e.Organizer.Name, Email: e.Organizer.Email}
	}
	for _, a := range e.Attendees {
		d.Attendees = append(d.Attendees, person{Name: a.Name, Email: a.Email})
	}
	return d
}

// validInterviewSource gates the persisted source tag. 'ai' covers both the
// free-text and screenshot parse paths (llmEventToDTO sets Source:"ai"), so it
// must be accepted here or AI-parsed events can be previewed but never saved.
func validInterviewSource(s string) bool {
	switch s {
	case "ics", "ai", "manual":
		return true
	}
	return false
}

type interviewCreateRequest struct {
	Source      string     `json:"source"`
	UID         string     `json:"uid"`
	Summary     string     `json:"summary"`
	Location    string     `json:"location"`
	Description string     `json:"description"`
	StartsAt    time.Time  `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
	AllDay      bool       `json:"all_day"`
	Organizer   *person    `json:"organizer,omitempty"`
	Attendees   []person   `json:"attendees"`
}

// POST /api/applications/{id}/interviews — persist a parsed (or manual) event.
// When source='ics' and uid is set, upserts on (application_id, uid).
func (s *Server) handleInterviewCreate(w http.ResponseWriter, r *http.Request) {
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

	var in interviewCreateRequest
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	in.Source = strings.TrimSpace(in.Source)
	if in.Source == "" {
		in.Source = "manual"
	}
	if !validInterviewSource(in.Source) {
		writeJSONError(w, http.StatusBadRequest, "source must be 'ics', 'ai', or 'manual'")
		return
	}
	in.Summary = strings.TrimSpace(in.Summary)
	if in.Summary == "" {
		in.Summary = "Interview"
	}
	if in.StartsAt.IsZero() {
		writeJSONError(w, http.StatusBadRequest, "starts_at is required")
		return
	}

	var uidPtr *string
	if u := strings.TrimSpace(in.UID); u != "" {
		uidPtr = &u
	}
	if in.Attendees == nil {
		in.Attendees = []person{}
	}
	attendeesJSON, _ := json.Marshal(in.Attendees)
	var organizerJSON []byte
	if in.Organizer != nil {
		organizerJSON, _ = json.Marshal(in.Organizer)
	}

	var locPtr, descPtr *string
	if loc := strings.TrimSpace(in.Location); loc != "" {
		locPtr = &loc
	}
	if desc := strings.TrimSpace(in.Description); desc != "" {
		descPtr = &desc
	}

	// Upsert on (application_id, uid) when uid is set. The partial unique index
	// makes the conflict target valid; without uid we always insert.
	var iv interviewDTO
	if uidPtr != nil {
		err = s.Pool.QueryRow(r.Context(), `
			INSERT INTO interviews (
			    application_id, user_id, source, uid, summary, location,
			    description, starts_at, ends_at, all_day, organizer, attendees
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
			ON CONFLICT (application_id, uid) WHERE uid IS NOT NULL DO UPDATE SET
			    source      = EXCLUDED.source,
			    summary     = EXCLUDED.summary,
			    location    = EXCLUDED.location,
			    description = EXCLUDED.description,
			    starts_at   = EXCLUDED.starts_at,
			    ends_at     = EXCLUDED.ends_at,
			    all_day     = EXCLUDED.all_day,
			    organizer   = EXCLUDED.organizer,
			    attendees   = EXCLUDED.attendees,
			    updated_at  = now()
			RETURNING id, application_id, source, uid, summary, location, description,
			    starts_at, ends_at, all_day, organizer, attendees, created_at`,
			appID, u.ID, in.Source, uidPtr, in.Summary, locPtr,
			descPtr, in.StartsAt, in.EndsAt, in.AllDay, organizerJSON, attendeesJSON,
		).Scan(&iv.ID, &iv.ApplicationID, &iv.Source, &iv.UID, &iv.Summary, &iv.Location,
			&iv.Description, &iv.StartsAt, &iv.EndsAt, &iv.AllDay, &iv.Organizer, &iv.Attendees,
			&iv.CreatedAt)
	} else {
		err = s.Pool.QueryRow(r.Context(), `
			INSERT INTO interviews (
			    application_id, user_id, source, uid, summary, location,
			    description, starts_at, ends_at, all_day, organizer, attendees
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
			RETURNING id, application_id, source, uid, summary, location, description,
			    starts_at, ends_at, all_day, organizer, attendees, created_at`,
			appID, u.ID, in.Source, uidPtr, in.Summary, locPtr,
			descPtr, in.StartsAt, in.EndsAt, in.AllDay, organizerJSON, attendeesJSON,
		).Scan(&iv.ID, &iv.ApplicationID, &iv.Source, &iv.UID, &iv.Summary, &iv.Location,
			&iv.Description, &iv.StartsAt, &iv.EndsAt, &iv.AllDay, &iv.Organizer, &iv.Attendees,
			&iv.CreatedAt)
	}
	if err != nil {
		s.Logger.Error("interview insert", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusCreated, iv)
}

// GET /api/applications/{id}/interviews — list all interviews for the app,
// future first then past, both in chronological order.
func (s *Server) handleInterviewsList(w http.ResponseWriter, r *http.Request) {
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
		SELECT id, application_id, source, uid, summary, location, description,
		       starts_at, ends_at, all_day, organizer, attendees, created_at
		FROM interviews
		WHERE application_id = $1 AND user_id = $2
		ORDER BY starts_at ASC`, appID, u.ID)
	if err != nil {
		s.Logger.Error("interviews list", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()

	out := []interviewDTO{}
	for rows.Next() {
		var iv interviewDTO
		if err := rows.Scan(&iv.ID, &iv.ApplicationID, &iv.Source, &iv.UID, &iv.Summary,
			&iv.Location, &iv.Description, &iv.StartsAt, &iv.EndsAt, &iv.AllDay,
			&iv.Organizer, &iv.Attendees, &iv.CreatedAt); err != nil {
			s.Logger.Error("interview scan", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		out = append(out, iv)
	}
	writeJSON(w, http.StatusOK, out)
}

// DELETE /api/applications/{id}/interviews/{iid}
func (s *Server) handleInterviewDelete(w http.ResponseWriter, r *http.Request) {
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
	tag, err := s.Pool.Exec(r.Context(),
		`DELETE FROM interviews WHERE id = $1 AND application_id = $2 AND user_id = $3`,
		iid, appID, u.ID)
	if err != nil {
		s.Logger.Error("interview delete", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if tag.RowsAffected() == 0 {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// nextInterview returns the soonest upcoming interview for an application
// (starts_at >= now). Returns (nil, nil) if none, so callers can branch on
// "do we have a real meeting" without an error check.
func (s *Server) nextInterview(ctx context.Context, appID int64) (*interviewDTO, error) {
	var iv interviewDTO
	err := s.Pool.QueryRow(ctx, `
		SELECT id, application_id, source, uid, summary, location, description,
		       starts_at, ends_at, all_day, organizer, attendees, created_at
		FROM interviews
		WHERE application_id = $1 AND starts_at >= now()
		ORDER BY starts_at ASC
		LIMIT 1`, appID,
	).Scan(&iv.ID, &iv.ApplicationID, &iv.Source, &iv.UID, &iv.Summary,
		&iv.Location, &iv.Description, &iv.StartsAt, &iv.EndsAt, &iv.AllDay,
		&iv.Organizer, &iv.Attendees, &iv.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &iv, nil
}
