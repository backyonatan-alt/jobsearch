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
	AppCount       int        `json:"app_count"`       // real applications, excludes [demo] seed rows
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
