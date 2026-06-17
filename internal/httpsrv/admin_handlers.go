package httpsrv

import (
	"net/http"
	"strconv"
	"time"
)

type inviteDTO struct {
	Email          string    `json:"email"`
	InvitedAt      time.Time `json:"invited_at"`
	Note           string    `json:"note,omitempty"`
	InvitedByEmail string    `json:"invited_by_email,omitempty"`
}

func (s *Server) handleAdminInvitesList(w http.ResponseWriter, r *http.Request) {
	invites, err := s.Auth.ListInvites(r.Context())
	if err != nil {
		s.Logger.Error("list invites", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	out := make([]inviteDTO, 0, len(invites))
	for _, i := range invites {
		d := inviteDTO{Email: i.Email, InvitedAt: i.InvitedAt, InvitedByEmail: i.InvitedByEmail}
		if i.Note != nil {
			d.Note = *i.Note
		}
		out = append(out, d)
	}
	writeJSON(w, http.StatusOK, out)
}

type addInviteRequest struct {
	Email string `json:"email"`
	Note  string `json:"note"`
}

func (s *Server) handleAdminInvitesAdd(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	var req addInviteRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.Auth.AddInvite(r.Context(), req.Email, req.Note, u.ID); err != nil {
		s.Logger.Info("add invite rejected", "err", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.Logger.Info("invite added", "by_user_id", u.ID, "email", req.Email)
	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

// ── Users + AI interview-prep credits ──────────────────────────────────────

type adminUserDTO struct {
	ID             int64      `json:"id"`
	Email          string     `json:"email"`
	IsAdmin        bool       `json:"is_admin"`
	PrepUsed       int        `json:"prep_credits_used"`
	PrepLimit      int        `json:"prep_credits_limit"`
	OnboardedAt    *time.Time `json:"onboarded_at"`
	CreatedAt      time.Time  `json:"created_at"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	AppCount       int        `json:"app_count"` // real applications, excludes [demo] seed rows
	InterviewCount int        `json:"interview_count"`
	DossierCount   int        `json:"dossier_count"`
}

func (s *Server) handleAdminUsersList(w http.ResponseWriter, r *http.Request) {
	// Per-user pilot-usage snapshot: when they joined, when they were last
	// seen, and what they've actually done in the product. Demo-seed rows
	// (notes prefixed "[demo] ") are excluded so the app count reflects real use.
	// Ordered by last sign-in so the most active pilots float to the top.
	rows, err := s.Pool.Query(r.Context(),
		`SELECT u.id, u.email, u.is_admin, u.prep_credits_used, u.prep_credits_limit,
		        u.onboarded_at, u.created_at, u.last_login_at,
		        (SELECT count(*) FROM applications a
		           WHERE a.user_id = u.id
		             AND (a.notes IS NULL OR a.notes NOT LIKE '[demo] %')) AS app_count,
		        (SELECT count(*) FROM interviews i WHERE i.user_id = u.id) AS interview_count,
		        (SELECT count(*) FROM dossiers d
		           JOIN applications a2 ON a2.id = d.application_id
		           WHERE a2.user_id = u.id) AS dossier_count
		 FROM users u
		 ORDER BY u.last_login_at DESC NULLS LAST, u.id`)
	if err != nil {
		s.Logger.Error("list users", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()
	out := make([]adminUserDTO, 0)
	for rows.Next() {
		var d adminUserDTO
		if err := rows.Scan(&d.ID, &d.Email, &d.IsAdmin, &d.PrepUsed, &d.PrepLimit, &d.OnboardedAt,
			&d.CreatedAt, &d.LastLoginAt, &d.AppCount, &d.InterviewCount, &d.DossierCount); err != nil {
			s.Logger.Error("scan user", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		out = append(out, d)
	}
	writeJSON(w, http.StatusOK, out)
}

// ── Invite funnel ──────────────────────────────────────────────────────────

type inviteFunnelDTO struct {
	Email          string     `json:"email"`
	Note           string     `json:"note,omitempty"`
	InvitedAt      time.Time  `json:"invited_at"`
	InvitedByEmail string     `json:"invited_by_email,omitempty"`
	SignedInAt     *time.Time `json:"signed_in_at,omitempty"` // users.created_at — first OAuth login
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	OnboardedAt    *time.Time `json:"onboarded_at,omitempty"`
	AppCount       int        `json:"app_count"` // real apps, excludes [demo] seed rows
	InterviewCount int        `json:"interview_count"`
	DossierCount   int        `json:"dossier_count"`
	EventCount     int        `json:"event_count"`
	LastActivityAt *time.Time `json:"last_activity_at,omitempty"` // max(last_login, last event)
	Stage          string     `json:"stage"`                      // invited | signed_in | activated | active | dormant
	PendingCount   int        `json:"pending_count,omitempty"`    // only on summary head
}

// handleAdminInviteFunnel returns one row per invited email joined to the user
// they became (matched on lowercased email), with the activity that tells you
// whether they actually use the product. Stage is the funnel position:
//
//	invited   — on the list, never signed in
//	signed_in — signed in, but created nothing yet
//	activated — created ≥1 app/interview/dossier
//	active    — activated AND seen in the last 7 days
//	dormant   — activated but not seen in 21+ days
func (s *Server) handleAdminInviteFunnel(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Pool.Query(r.Context(), `
		SELECT i.email, i.note, i.invited_at, inv.email AS invited_by_email,
		       u.created_at AS signed_in_at, u.last_login_at, u.onboarded_at,
		       COALESCE(app.c, 0)  AS app_count,
		       COALESCE(iv.c, 0)   AS interview_count,
		       COALESCE(dos.c, 0)  AS dossier_count,
		       COALESCE(ev.c, 0)   AS event_count,
		       ev.last_at          AS last_event_at
		FROM invited_emails i
		LEFT JOIN users inv ON inv.id = i.invited_by_user_id
		LEFT JOIN users u   ON lower(u.email) = lower(i.email)
		LEFT JOIN LATERAL (SELECT count(*) c FROM applications a
		    WHERE a.user_id = u.id
		      AND (a.notes IS NULL OR a.notes NOT LIKE '[demo] %')) app ON true
		LEFT JOIN LATERAL (SELECT count(*) c FROM interviews x WHERE x.user_id = u.id) iv ON true
		LEFT JOIN LATERAL (SELECT count(*) c FROM dossiers d
		    JOIN applications a2 ON a2.id = d.application_id
		    WHERE a2.user_id = u.id) dos ON true
		LEFT JOIN LATERAL (SELECT count(*) c, max(created_at) last_at
		    FROM events e WHERE e.user_id = u.id) ev ON true
		ORDER BY i.invited_at DESC`)
	if err != nil {
		s.Logger.Error("invite funnel", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()

	now := time.Now()
	out := make([]inviteFunnelDTO, 0)
	for rows.Next() {
		var d inviteFunnelDTO
		var note, invitedBy *string
		var lastEventAt *time.Time
		if err := rows.Scan(&d.Email, &note, &d.InvitedAt, &invitedBy,
			&d.SignedInAt, &d.LastLoginAt, &d.OnboardedAt,
			&d.AppCount, &d.InterviewCount, &d.DossierCount, &d.EventCount, &lastEventAt); err != nil {
			s.Logger.Error("invite funnel scan", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		if note != nil {
			d.Note = *note
		}
		if invitedBy != nil {
			d.InvitedByEmail = *invitedBy
		}
		// Last activity = most recent of last login and last logged event.
		d.LastActivityAt = d.LastLoginAt
		if lastEventAt != nil && (d.LastActivityAt == nil || lastEventAt.After(*d.LastActivityAt)) {
			d.LastActivityAt = lastEventAt
		}
		d.Stage = inviteStage(d, now)
		out = append(out, d)
	}

	// Top-of-funnel: people who asked for access but aren't invited yet.
	var pending int
	_ = s.Pool.QueryRow(r.Context(),
		`SELECT count(*) FROM beta_interest WHERE invited_at IS NULL`).Scan(&pending)

	writeJSON(w, http.StatusOK, map[string]any{
		"pending_count": pending,
		"invitees":      out,
	})
}

func inviteStage(d inviteFunnelDTO, now time.Time) string {
	if d.SignedInAt == nil {
		return "invited"
	}
	activated := d.AppCount > 0 || d.InterviewCount > 0 || d.DossierCount > 0
	if !activated {
		return "signed_in"
	}
	if d.LastActivityAt != nil {
		if now.Sub(*d.LastActivityAt) <= 7*24*time.Hour {
			return "active"
		}
		if now.Sub(*d.LastActivityAt) >= 21*24*time.Hour {
			return "dormant"
		}
	}
	return "activated"
}

type grantPrepRequest struct {
	Add int `json:"add"`
}

func (s *Server) handleAdminGrantPrep(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	var req grantPrepRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Add < 1 || req.Add > 1000 {
		writeJSONError(w, http.StatusBadRequest, "add must be between 1 and 1000")
		return
	}
	var used, limit int
	err = s.Pool.QueryRow(r.Context(),
		`UPDATE users SET prep_credits_limit = prep_credits_limit + $2 WHERE id = $1
		 RETURNING prep_credits_used, prep_credits_limit`, id, req.Add).Scan(&used, &limit)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "user not found")
		return
	}
	s.Logger.Info("prep credits granted", "user_id", id, "add", req.Add, "new_limit", limit)
	writeJSON(w, http.StatusOK, map[string]int{"prep_credits_used": used, "prep_credits_limit": limit})
}

// ── Adoption / activation ──────────────────────────────────────────────────

type milestoneDTO struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Users int    `json:"users"` // distinct users who reached this milestone
}

type eventStatDTO struct {
	Name   string     `json:"name"`
	Users  int        `json:"users"`  // distinct users who fired it
	Total  int        `json:"total"`  // all-time count
	Recent int        `json:"recent"` // count in the last 7 days
	LastAt *time.Time `json:"last_at"`
}

type adoptionResp struct {
	TotalUsers int            `json:"total_users"`
	Milestones []milestoneDTO `json:"milestones"`
	Events     []eventStatDTO `json:"events"`
}

// handleAdminAdoption answers "of everyone who signed in, how many reached each
// AI-wedge moment, and which surfaces are actually used." Milestones are counted
// from the real tables (demo rows excluded) so they reflect genuine activation;
// the events table gives the raw per-surface usage, including dead surfaces.
func (s *Server) handleAdminAdoption(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var resp adoptionResp

	var total, createdApp, addedInterview, genDossier, ranPrep, bulkImport int
	err := s.Pool.QueryRow(ctx, `
		SELECT
		  (SELECT count(*) FROM users),
		  (SELECT count(DISTINCT user_id) FROM applications
		     WHERE notes IS NULL OR notes NOT LIKE '[demo] %'),
		  (SELECT count(DISTINCT user_id) FROM interviews),
		  (SELECT count(DISTINCT a.user_id) FROM dossiers d
		     JOIN applications a ON a.id = d.application_id),
		  (SELECT count(*) FROM users WHERE prep_credits_used > 0),
		  (SELECT count(DISTINCT user_id) FROM events WHERE name = 'bulk_import')
	`).Scan(&total, &createdApp, &addedInterview, &genDossier, &ranPrep, &bulkImport)
	if err != nil {
		s.Logger.Error("adoption milestones", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	resp.TotalUsers = total
	resp.Milestones = []milestoneDTO{
		{Key: "created_app", Label: "Created an application", Users: createdApp},
		{Key: "added_interview", Label: "Added an interview", Users: addedInterview},
		{Key: "generated_dossier", Label: "Generated a dossier", Users: genDossier},
		{Key: "ran_prep", Label: "Ran interview prep", Users: ranPrep},
		{Key: "bulk_import", Label: "Used bulk import", Users: bulkImport},
	}

	rows, err := s.Pool.Query(ctx, `
		SELECT name, count(*) AS total, count(DISTINCT user_id) AS users,
		       count(*) FILTER (WHERE created_at > now() - interval '7 days') AS recent,
		       max(created_at) AS last_at
		FROM events
		GROUP BY name
		ORDER BY users DESC, total DESC`)
	if err != nil {
		s.Logger.Error("adoption events", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()
	resp.Events = make([]eventStatDTO, 0)
	for rows.Next() {
		var e eventStatDTO
		if err := rows.Scan(&e.Name, &e.Total, &e.Users, &e.Recent, &e.LastAt); err != nil {
			s.Logger.Error("adoption events scan", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		resp.Events = append(resp.Events, e)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleAdminInvitesDelete(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	if email == "" {
		writeJSONError(w, http.StatusBadRequest, "missing email")
		return
	}
	if err := s.Auth.RemoveInvite(r.Context(), email); err != nil {
		s.Logger.Error("remove invite", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
