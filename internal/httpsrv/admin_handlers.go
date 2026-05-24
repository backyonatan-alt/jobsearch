package httpsrv

import (
	"net/http"
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
