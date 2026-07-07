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

// dossierRow is one persisted brief — the shared company brief when InterviewID
// is nil, an interviewer brief for a specific round when set.
type dossierRow struct {
	InterviewerName *string
	Content         json.RawMessage
	ModelUsed       string
	GeneratedAt     time.Time
	InterviewID     *int64
}

const dossierCols = `interviewer_name, content, model_used, generated_at, interview_id`

func scanDossier(row interface{ Scan(...any) error }) (*dossierRow, error) {
	var d dossierRow
	err := row.Scan(&d.InterviewerName, &d.Content, &d.ModelUsed, &d.GeneratedAt, &d.InterviewID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// fetchCompanyBrief loads the application-level company brief (interview_id NULL).
func (s *Server) fetchCompanyBrief(ctx context.Context, appID int64) (*dossierRow, error) {
	return scanDossier(s.Pool.QueryRow(ctx,
		`SELECT `+dossierCols+` FROM dossiers WHERE application_id = $1 AND interview_id IS NULL`, appID))
}

// fetchRoundBrief loads the interviewer brief for one round.
func (s *Server) fetchRoundBrief(ctx context.Context, appID, interviewID int64) (*dossierRow, error) {
	return scanDossier(s.Pool.QueryRow(ctx,
		`SELECT `+dossierCols+` FROM dossiers WHERE application_id = $1 AND interview_id = $2`, appID, interviewID))
}

// composeContent merges the shared company block into an interviewer brief so
// the Today dashboard can read snapshot/lands AND company.watch_fors from a
// single content object. Either side may be nil.
func composeContent(interviewer, company json.RawMessage) json.RawMessage {
	if len(company) == 0 {
		return interviewer
	}
	if len(interviewer) == 0 {
		return company
	}
	var im map[string]json.RawMessage
	var cm map[string]json.RawMessage
	if json.Unmarshal(interviewer, &im) != nil || json.Unmarshal(company, &cm) != nil {
		return interviewer
	}
	if cb, ok := cm["company"]; ok {
		im["company"] = cb
	}
	merged, err := json.Marshal(im)
	if err != nil {
		return interviewer
	}
	return merged
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

	q := r.URL.Query()

	// ?scope=company → the shared company brief tab (no fallback).
	if q.Get("scope") == "company" {
		d, err := s.fetchCompanyBrief(r.Context(), appID)
		if err != nil {
			s.Logger.Error("dossier get", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		if d == nil {
			writeJSONError(w, http.StatusNotFound, "no company brief yet")
			return
		}
		writeDossier(w, appID, d, s.meetingForApp(r.Context(), app, ""))
		return
	}

	// ?interview_id=N → one round's interviewer brief (no fallback — empty if none).
	if qid := q.Get("interview_id"); qid != "" {
		iid, perr := strconv.ParseInt(qid, 10, 64)
		if perr != nil {
			writeJSONError(w, http.StatusBadRequest, "bad interview_id")
			return
		}
		d, err := s.fetchRoundBrief(r.Context(), appID, iid)
		if err != nil {
			s.Logger.Error("dossier get", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		if d == nil {
			writeJSONError(w, http.StatusNotFound, "no prep for this round yet")
			return
		}
		writeDossier(w, appID, d, s.meetingForDossier(r.Context(), app, d.InterviewID, deref(d.InterviewerName)))
		return
	}

	// No params (Today dashboard / initial load): compose the next round's
	// interviewer brief with the shared company block so one content object has
	// both snapshot/lands and company.watch_fors.
	var nextID *int64
	if iv, e := s.nextInterview(r.Context(), appID); e == nil && iv != nil {
		nextID = &iv.ID
	}
	company, err := s.fetchCompanyBrief(r.Context(), appID)
	if err != nil {
		s.Logger.Error("dossier get", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	var round *dossierRow
	if nextID != nil {
		if round, err = s.fetchRoundBrief(r.Context(), appID, *nextID); err != nil {
			s.Logger.Error("dossier get", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
	}
	if round == nil && company == nil {
		writeJSONError(w, http.StatusNotFound, "no dossier yet")
		return
	}
	base := round
	if base == nil {
		base = company
	}
	var companyContent json.RawMessage
	if company != nil {
		companyContent = company.Content
	}
	writeJSON(w, http.StatusOK, dossierDTO{
		ApplicationID:   appID,
		InterviewID:     base.InterviewID,
		InterviewerName: deref(base.InterviewerName),
		GeneratedAt:     base.GeneratedAt,
		GeneratedAgo:    humanAgo(base.GeneratedAt),
		ModelUsed:       base.ModelUsed,
		Content:         composeContent(round.contentOrNil(), companyContent),
		Meeting:         s.meetingForDossier(r.Context(), app, base.InterviewID, deref(base.InterviewerName)),
	})
}

func (d *dossierRow) contentOrNil() json.RawMessage {
	if d == nil {
		return nil
	}
	return d.Content
}

func writeDossier(w http.ResponseWriter, appID int64, d *dossierRow, meeting meetingDTO) {
	writeJSON(w, http.StatusOK, dossierDTO{
		ApplicationID:   appID,
		InterviewID:     d.InterviewID,
		InterviewerName: deref(d.InterviewerName),
		GeneratedAt:     d.GeneratedAt,
		GeneratedAgo:    humanAgo(d.GeneratedAt),
		ModelUsed:       d.ModelUsed,
		Content:         d.Content,
		Meeting:         meeting,
	})
}

type refreshRequest struct {
	InterviewerName string `json:"interviewer_name"`
	InterviewID     *int64 `json:"interview_id"`
	// CompanyURL is an optional user-confirmed company website ("Not them? →"),
	// the authoritative grounding signal that overrides same-named-company drift.
	CompanyURL string `json:"company_url"`
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

	// If a round was named, it must belong to this application. Keep the row —
	// we need its starts_at to pull earlier rounds' debriefs (feed-forward).
	var round *interviewDTO
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
		round = iv
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
		"interviewer", req.InterviewerName, "round", req.InterviewID)

	now := time.Now()
	const modelUsed = "claude-sonnet-4-6+web_search"

	// Grounding signals so research pins the RIGHT same-named company.
	loc, jd := deref(app.Location), deref(app.JDURL)

	// Generation runs on a context detached from the connection: a dropped
	// client (phone locks, tab closes, network blip — the browser's "Failed to
	// fetch") must not abort a 1–2 min LLM call that's already paid for. The
	// brief is persisted regardless; the frontend recovers by polling GET.
	gctx := context.WithoutCancel(r.Context())

	// No round → (re)generate the shared company brief.
	if req.InterviewID == nil {
		content, gerr := s.LLM.GenerateCompanyBrief(gctx, app.Company, app.Role, loc, jd, req.CompanyURL)
		if gerr != nil {
			s.failGenerate(w, gctx, u.ID, start, gerr)
			return
		}
		if err := s.upsertCompanyBrief(gctx, appID, content, modelUsed, now); err != nil {
			s.Logger.Error("company brief upsert", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		s.finishGenerate(gctx, u.ID, start, len(content))
		writeDossier(w, appID, &dossierRow{Content: content, ModelUsed: modelUsed, GeneratedAt: now},
			s.meetingForApp(gctx, app, ""))
		return
	}

	// A round → research this interviewer. If the application has no company
	// brief yet, generate it alongside (in parallel) — bundled, no extra credit.
	var companyCh chan json.RawMessage
	if existing, _ := s.fetchCompanyBrief(gctx, appID); existing == nil {
		companyCh = make(chan json.RawMessage, 1)
		go func() {
			cc, cerr := s.LLM.GenerateCompanyBrief(gctx, app.Company, app.Role, loc, jd, req.CompanyURL)
			if cerr != nil {
				s.Logger.Info("bundled company brief failed (non-fatal)", "err", cerr)
				companyCh <- nil
				return
			}
			companyCh <- cc
		}()
	}

	// Feed-forward: fold in debriefs from rounds that happened before this one.
	var priorDebriefs string
	if round != nil {
		priorDebriefs = s.priorDebriefsContext(gctx, appID, round.StartsAt)
	}

	content, gerr := s.LLM.GenerateInterviewerBrief(gctx, app.Company, app.Role, req.InterviewerName, loc, req.CompanyURL, priorDebriefs)
	if gerr != nil {
		s.failGenerate(w, gctx, u.ID, start, gerr)
		return
	}
	var nullableName *string
	if req.InterviewerName != "" {
		nullableName = &req.InterviewerName
	}
	if err := s.upsertRoundBrief(gctx, appID, *req.InterviewID, nullableName, content, modelUsed, now); err != nil {
		s.Logger.Error("round brief upsert", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if companyCh != nil {
		if cc := <-companyCh; len(cc) > 0 {
			if uerr := s.upsertCompanyBrief(gctx, appID, cc, modelUsed, now); uerr != nil {
				s.Logger.Error("bundled company brief upsert", "err", uerr)
			}
		}
	}
	s.finishGenerate(gctx, u.ID, start, len(content))
	writeDossier(w, appID,
		&dossierRow{InterviewerName: nullableName, Content: content, ModelUsed: modelUsed, GeneratedAt: now, InterviewID: req.InterviewID},
		s.meetingForDossier(gctx, app, req.InterviewID, req.InterviewerName))
}

func (s *Server) upsertCompanyBrief(ctx context.Context, appID int64, content json.RawMessage, model string, now time.Time) error {
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO dossiers (application_id, interview_id, interviewer_name, content, model_used, generated_at)
		VALUES ($1, NULL, NULL, $2, $3, $4)
		ON CONFLICT (application_id) WHERE interview_id IS NULL DO UPDATE SET
		    content      = EXCLUDED.content,
		    model_used   = EXCLUDED.model_used,
		    generated_at = EXCLUDED.generated_at`,
		appID, content, model, now)
	return err
}

func (s *Server) upsertRoundBrief(ctx context.Context, appID, interviewID int64, name *string, content json.RawMessage, model string, now time.Time) error {
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO dossiers (application_id, interview_id, interviewer_name, content, model_used, generated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (interview_id) WHERE interview_id IS NOT NULL DO UPDATE SET
		    interviewer_name = EXCLUDED.interviewer_name,
		    content          = EXCLUDED.content,
		    model_used       = EXCLUDED.model_used,
		    generated_at     = EXCLUDED.generated_at`,
		appID, interviewID, name, content, model, now)
	return err
}

func (s *Server) failGenerate(w http.ResponseWriter, ctx context.Context, userID int64, start time.Time, err error) {
	s.Logger.Info("dossier generate failed", "err", err)
	s.logEvent(ctx, userID, "dossier_refresh", map[string]any{
		"outcome": "error", "error_reason": "generate_failed",
		"duration_ms": time.Since(start).Milliseconds(),
	})
	writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
}

// finishGenerate logs success and counts one prep credit. A bundled company
// brief rides along free — one user action, one credit.
func (s *Server) finishGenerate(ctx context.Context, userID int64, start time.Time, contentBytes int) {
	s.Logger.Info("dossier generate done", "bytes", contentBytes)
	s.logEvent(ctx, userID, "dossier_refresh", map[string]any{
		"outcome": "success", "duration_ms": time.Since(start).Milliseconds(),
	})
	if _, err := s.Pool.Exec(ctx,
		`UPDATE users SET prep_credits_used = prep_credits_used + 1 WHERE id = $1`, userID,
	); err != nil {
		s.Logger.Error("prep credits increment", "err", err) // non-fatal
	}
}

// applicationRow is the minimal application row the dossier endpoints need.
type applicationRow struct {
	ID       int64
	Company  string
	Role     string
	Status   string
	Location *string
	JDURL    *string
}

func (s *Server) fetchApplication(ctx context.Context, userID, appID int64) (*applicationRow, error) {
	var a applicationRow
	err := s.Pool.QueryRow(ctx, `
		SELECT id, company, role, status, location, jd_url
		FROM applications WHERE id = $1 AND user_id = $2`, appID, userID,
	).Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.Location, &a.JDURL)
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
