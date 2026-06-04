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
	ID          int64      `json:"id"`
	Email       string     `json:"email"`
	IsAdmin     bool       `json:"is_admin"`
	PrepUsed    int        `json:"prep_credits_used"`
	PrepLimit   int        `json:"prep_credits_limit"`
	OnboardedAt *time.Time `json:"onboarded_at"`
}

func (s *Server) handleAdminUsersList(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Pool.Query(r.Context(),
		`SELECT id, email, is_admin, prep_credits_used, prep_credits_limit, onboarded_at
		 FROM users ORDER BY id`)
	if err != nil {
		s.Logger.Error("list users", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()
	out := make([]adminUserDTO, 0)
	for rows.Next() {
		var d adminUserDTO
		if err := rows.Scan(&d.ID, &d.Email, &d.IsAdmin, &d.PrepUsed, &d.PrepLimit, &d.OnboardedAt); err != nil {
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
