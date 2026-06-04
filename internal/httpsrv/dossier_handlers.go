package httpsrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// Dossier as returned by the HTTP API. The `content` field is whatever Claude
// produced (interviewer / snapshot / signals / style / lands / avoid /
// questions); we mix in meeting + generated_at + interviewer_name on top.
type dossierDTO struct {
	ApplicationID   int64           `json:"application_id"`
	InterviewerName string          `json:"interviewer_name,omitempty"`
	GeneratedAt     time.Time       `json:"generated_at"`
	GeneratedAgo    string          `json:"generatedAgo"`
	ModelUsed       string          `json:"model_used"`
	Content         json.RawMessage `json:"content"`
	Meeting         meetingDTO      `json:"meeting"`
}

type meetingDTO struct {
	// When/Duration are server-rendered fallbacks for when no real interview
	// is linked — placeholder strings the dossier hero used before .ics
	// ingest existed. When StartsAt is set the frontend formats from the raw
	// timestamps instead so the viewer sees their own timezone.
	When     string `json:"when"`
	Duration string `json:"duration"`
	Medium   string `json:"medium"`
	Panel    string `json:"panel"`

	StartsAt *time.Time `json:"starts_at,omitempty"`
	EndsAt   *time.Time `json:"ends_at,omitempty"`
	AllDay   bool       `json:"all_day,omitempty"`
}

// GET /api/applications/:id/dossier — returns the cached dossier or 404.
func (s *Server) handleDossierGet(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	appID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	app, err := s.fetchApplication(r.Context(), u.ID, appID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}

	var (
		interviewerName *string
		content         json.RawMessage
		modelUsed       string
		generatedAt     time.Time
	)
	err = s.Pool.QueryRow(r.Context(), `
		SELECT interviewer_name, content, model_used, generated_at
		FROM dossiers WHERE application_id = $1`, appID,
	).Scan(&interviewerName, &content, &modelUsed, &generatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSONError(w, http.StatusNotFound, "no dossier yet")
		return
	}
	if err != nil {
		s.Logger.Error("dossier get", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}

	writeJSON(w, http.StatusOK, dossierDTO{
		ApplicationID:   appID,
		InterviewerName: deref(interviewerName),
		GeneratedAt:     generatedAt,
		GeneratedAgo:    humanAgo(generatedAt),
		ModelUsed:       modelUsed,
		Content:         content,
		Meeting:         s.meetingForApp(r.Context(), app, deref(interviewerName)),
	})
}

type refreshRequest struct {
	InterviewerName string `json:"interviewer_name"`
}

// POST /api/applications/:id/dossier/refresh — generate fresh, upsert, return.
func (s *Server) handleDossierRefresh(w http.ResponseWriter, r *http.Request) {
	if s.LLM == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "AI dossier generation is not configured (ANTHROPIC_API_KEY missing)")
		return
	}
	u, _ := userFromCtx(r.Context())
	appID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	app, err := s.fetchApplication(r.Context(), u.ID, appID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}

	var req refreshRequest
	if r.ContentLength > 0 {
		if err := readJSON(r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	req.InterviewerName = strings.TrimSpace(req.InterviewerName)

	// Per-user cap on AI generations (beta cost control). Admin can grant more.
	var prepUsed, prepLimit int
	if err := s.Pool.QueryRow(r.Context(),
		`SELECT prep_credits_used, prep_credits_limit FROM users WHERE id = $1`, u.ID,
	).Scan(&prepUsed, &prepLimit); err != nil {
		s.Logger.Error("prep credits read", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if prepUsed >= prepLimit {
		writeJSONError(w, http.StatusTooManyRequests,
			"You've reached your interview-prep limit for the beta. Ask the admin to add more.")
		return
	}

	s.Logger.Info("dossier generate start",
		"user_id", u.ID, "app_id", appID,
		"company", app.Company, "role", app.Role,
		"interviewer", req.InterviewerName)

	content, err := s.LLM.GenerateDossier(r.Context(), app.Company, app.Role, req.InterviewerName)
	if err != nil {
		s.Logger.Info("dossier generate failed", "err", err)
		writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	now := time.Now()
	const modelUsed = "claude-sonnet-4-6+web_search"
	var nullableName *string
	if req.InterviewerName != "" {
		nullableName = &req.InterviewerName
	}
	_, err = s.Pool.Exec(r.Context(), `
		INSERT INTO dossiers (application_id, interviewer_name, content, model_used, generated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (application_id) DO UPDATE SET
		    interviewer_name = EXCLUDED.interviewer_name,
		    content          = EXCLUDED.content,
		    model_used       = EXCLUDED.model_used,
		    generated_at     = EXCLUDED.generated_at`,
		appID, nullableName, content, modelUsed, now)
	if err != nil {
		s.Logger.Error("dossier upsert", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	s.Logger.Info("dossier generate done", "app_id", appID, "bytes", len(content))

	// Count this generation against the user's beta cap.
	if _, err := s.Pool.Exec(r.Context(),
		`UPDATE users SET prep_credits_used = prep_credits_used + 1 WHERE id = $1`, u.ID,
	); err != nil {
		s.Logger.Error("prep credits increment", "err", err) // non-fatal
	}

	writeJSON(w, http.StatusOK, dossierDTO{
		ApplicationID:   appID,
		InterviewerName: req.InterviewerName,
		GeneratedAt:     now,
		GeneratedAgo:    "just now",
		ModelUsed:       modelUsed,
		Content:         content,
		Meeting:         s.meetingForApp(r.Context(), app, req.InterviewerName),
	})
}

// applicationRow is the minimal application row the dossier endpoints need.
type applicationRow struct {
	ID      int64
	Company string
	Role    string
	Status  string
}

func (s *Server) fetchApplication(ctx context.Context, userID, appID int64) (*applicationRow, error) {
	var a applicationRow
	err := s.Pool.QueryRow(ctx, `
		SELECT id, company, role, status
		FROM applications WHERE id = $1 AND user_id = $2`, appID, userID,
	).Scan(&a.ID, &a.Company, &a.Role, &a.Status)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func meetingPlaceholder(app *applicationRow, interviewer string) meetingDTO {
	panel := interviewer
	if panel == "" {
		panel = "(hiring team)"
	}
	return meetingDTO{
		When:     "Upcoming · time TBD",
		Duration: "60 min",
		Medium:   "Video call",
		Panel:    panel,
	}
}

// meetingForApp prefers a real upcoming interview pulled from the interviews
// table; if none exists, falls back to the placeholder so the dossier still
// renders. Logged-and-swallowed errors keep the dossier resilient even if the
// interviews query fails.
//
// When a real interview exists, we emit only the raw timestamps (StartsAt /
// EndsAt / AllDay) and leave When/Duration empty — the frontend formats both
// in the viewer's timezone. Server-side formatting renders in the server's
// zone, which surfaces as a TZ mismatch against the Scheduled list.
func (s *Server) meetingForApp(ctx context.Context, app *applicationRow, interviewer string) meetingDTO {
	iv, err := s.nextInterview(ctx, app.ID)
	if err != nil {
		s.Logger.Error("next interview lookup", "err", err, "app_id", app.ID)
	}
	if iv == nil {
		return meetingPlaceholder(app, interviewer)
	}

	panel := interviewer
	if panel == "" {
		panel = "(hiring team)"
	}
	medium := "Video call"
	if loc := deref(iv.Location); loc != "" {
		medium = loc
	}
	start := iv.StartsAt
	return meetingDTO{
		Panel:    panel,
		Medium:   medium,
		StartsAt: &start,
		EndsAt:   iv.EndsAt,
		AllDay:   iv.AllDay,
	}
}

func humanAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
