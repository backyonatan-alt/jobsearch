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
	InterviewID     *int64          `json:"interview_id,omitempty"`
	InterviewerName string          `json:"interviewer_name,omitempty"`
	GeneratedAt     time.Time       `json:"generated_at"`
	GeneratedAgo    string          `json:"generatedAgo"`
	ModelUsed       string          `json:"model_used"`
	Content         json.RawMessage `json:"content"`
	Meeting         meetingDTO      `json:"meeting"`
}

// dossierRow is one persisted dossier (round-specific when InterviewID is set,
// application-level "General" prep when nil).
type dossierRow struct {
	InterviewerName *string
	Content         json.RawMessage
	ModelUsed       string
	GeneratedAt     time.Time
	InterviewID     *int64
}

// fetchDossier loads the dossier for a round. When interviewID is set it looks
// for that round's prep; if none and allowFallback is true it falls back to the
// application-level (General) prep. Returns (nil, nil) when nothing matches.
func (s *Server) fetchDossier(ctx context.Context, appID int64, interviewID *int64, allowFallback bool) (*dossierRow, error) {
	const cols = `interviewer_name, content, model_used, generated_at, interview_id`
	if interviewID != nil {
		var d dossierRow
		err := s.Pool.QueryRow(ctx,
			`SELECT `+cols+` FROM dossiers WHERE application_id = $1 AND interview_id = $2`,
			appID, *interviewID,
		).Scan(&d.InterviewerName, &d.Content, &d.ModelUsed, &d.GeneratedAt, &d.InterviewID)
		if err == nil {
			return &d, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		if !allowFallback {
			return nil, nil
		}
	}
	var d dossierRow
	err := s.Pool.QueryRow(ctx,
		`SELECT `+cols+` FROM dossiers WHERE application_id = $1 AND interview_id IS NULL`,
		appID,
	).Scan(&d.InterviewerName, &d.Content, &d.ModelUsed, &d.GeneratedAt, &d.InterviewID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
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

	// Which round's prep? An explicit ?interview_id pins one round (no fallback —
	// "show me this round, empty if none"). Without it we resolve to the next
	// upcoming round and fall back to the General prep.
	var interviewID *int64
	allowFallback := false
	if q := r.URL.Query().Get("interview_id"); q != "" {
		v, perr := strconv.ParseInt(q, 10, 64)
		if perr != nil {
			writeJSONError(w, http.StatusBadRequest, "bad interview_id")
			return
		}
		interviewID = &v
	} else {
		if iv, e := s.nextInterview(r.Context(), appID); e == nil && iv != nil {
			interviewID = &iv.ID
		}
		allowFallback = true
	}

	d, err := s.fetchDossier(r.Context(), appID, interviewID, allowFallback)
	if err != nil {
		s.Logger.Error("dossier get", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if d == nil {
		writeJSONError(w, http.StatusNotFound, "no dossier yet")
		return
	}

	writeJSON(w, http.StatusOK, dossierDTO{
		ApplicationID:   appID,
		InterviewID:     d.InterviewID,
		InterviewerName: deref(d.InterviewerName),
		GeneratedAt:     d.GeneratedAt,
		GeneratedAgo:    humanAgo(d.GeneratedAt),
		ModelUsed:       d.ModelUsed,
		Content:         d.Content,
		Meeting:         s.meetingForDossier(r.Context(), app, d.InterviewID, deref(d.InterviewerName)),
	})
}

type refreshRequest struct {
	InterviewerName string `json:"interviewer_name"`
	InterviewID     *int64 `json:"interview_id"`
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
	start := time.Now()

	var req refreshRequest
	if r.ContentLength > 0 {
		if err := readJSON(r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	req.InterviewerName = strings.TrimSpace(req.InterviewerName)

	// If a round was named, it must belong to this application.
	if req.InterviewID != nil {
		iv, ierr := s.interviewByID(r.Context(), appID, *req.InterviewID)
		if ierr != nil {
			s.Logger.Error("interview lookup", "err", ierr)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		if iv == nil {
			writeJSONError(w, http.StatusNotFound, "interview not found")
			return
		}
	}

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
		s.logEvent(r.Context(), u.ID, "dossier_refresh", map[string]any{
			"outcome": "limit", "duration_ms": time.Since(start).Milliseconds(),
		})
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
		s.logEvent(r.Context(), u.ID, "dossier_refresh", map[string]any{
			"outcome": "error", "error_reason": "generate_failed",
			"duration_ms": time.Since(start).Milliseconds(),
		})
		writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	now := time.Now()
	const modelUsed = "claude-sonnet-4-6+web_search"
	var nullableName *string
	if req.InterviewerName != "" {
		nullableName = &req.InterviewerName
	}
	if req.InterviewID != nil {
		_, err = s.Pool.Exec(r.Context(), `
			INSERT INTO dossiers (application_id, interview_id, interviewer_name, content, model_used, generated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (interview_id) WHERE interview_id IS NOT NULL DO UPDATE SET
			    interviewer_name = EXCLUDED.interviewer_name,
			    content          = EXCLUDED.content,
			    model_used       = EXCLUDED.model_used,
			    generated_at     = EXCLUDED.generated_at`,
			appID, *req.InterviewID, nullableName, content, modelUsed, now)
	} else {
		_, err = s.Pool.Exec(r.Context(), `
			INSERT INTO dossiers (application_id, interview_id, interviewer_name, content, model_used, generated_at)
			VALUES ($1, NULL, $2, $3, $4, $5)
			ON CONFLICT (application_id) WHERE interview_id IS NULL DO UPDATE SET
			    interviewer_name = EXCLUDED.interviewer_name,
			    content          = EXCLUDED.content,
			    model_used       = EXCLUDED.model_used,
			    generated_at     = EXCLUDED.generated_at`,
			appID, nullableName, content, modelUsed, now)
	}
	if err != nil {
		s.Logger.Error("dossier upsert", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	s.Logger.Info("dossier generate done", "app_id", appID, "bytes", len(content))
	s.logEvent(r.Context(), u.ID, "dossier_refresh", map[string]any{
		"outcome": "success", "duration_ms": time.Since(start).Milliseconds(),
	})

	// Count this generation against the user's beta cap.
	if _, err := s.Pool.Exec(r.Context(),
		`UPDATE users SET prep_credits_used = prep_credits_used + 1 WHERE id = $1`, u.ID,
	); err != nil {
		s.Logger.Error("prep credits increment", "err", err) // non-fatal
	}

	writeJSON(w, http.StatusOK, dossierDTO{
		ApplicationID:   appID,
		InterviewID:     req.InterviewID,
		InterviewerName: req.InterviewerName,
		GeneratedAt:     now,
		GeneratedAgo:    "just now",
		ModelUsed:       modelUsed,
		Content:         content,
		Meeting:         s.meetingForDossier(r.Context(), app, req.InterviewID, req.InterviewerName),
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
	return meetingFromInterview(iv, interviewer)
}

// meetingForDossier renders the meeting block for the round this dossier is
// pinned to. General (no interview_id) prep falls back to the next upcoming
// interview, preserving the pre-per-round behaviour.
func (s *Server) meetingForDossier(ctx context.Context, app *applicationRow, interviewID *int64, interviewer string) meetingDTO {
	if interviewID != nil {
		iv, err := s.interviewByID(ctx, app.ID, *interviewID)
		if err != nil {
			s.Logger.Error("interview lookup", "err", err, "app_id", app.ID)
		}
		if iv != nil {
			return meetingFromInterview(iv, interviewer)
		}
	}
	return s.meetingForApp(ctx, app, interviewer)
}

func meetingFromInterview(iv *interviewDTO, interviewer string) meetingDTO {
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
