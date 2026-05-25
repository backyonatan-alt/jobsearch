package httpsrv

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type betaInterestRequest struct {
	Email  string `json:"email"`
	Note   string `json:"note"`
	Source string `json:"source"`
}

// handleBetaInterestCreate is public — anyone can drop their email here.
// Idempotent on email: re-submissions update note/source rather than 409,
// because the user almost always wants to add context after the fact.
func (s *Server) handleBetaInterestCreate(w http.ResponseWriter, r *http.Request) {
	var req betaInterestRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	note := strings.TrimSpace(req.Note)
	source := strings.TrimSpace(req.Source)
	if !emailRe.MatchString(email) || len(email) > 320 {
		writeJSONError(w, http.StatusBadRequest, "enter a valid email")
		return
	}
	if len(note) > 500 {
		note = note[:500]
	}
	if len(source) > 80 {
		source = source[:80]
	}

	// If they're already invited, tell them so — no need for a pending row.
	var already bool
	if err := s.Pool.QueryRow(r.Context(),
		`SELECT EXISTS (SELECT 1 FROM invited_emails WHERE email = $1)`, email,
	).Scan(&already); err == nil && already {
		writeJSON(w, http.StatusOK, map[string]any{"status": "already_invited"})
		return
	}

	_, err := s.Pool.Exec(r.Context(), `
		INSERT INTO beta_interest (email, note, source)
		VALUES ($1, NULLIF($2,''), NULLIF($3,''))
		ON CONFLICT (email) DO UPDATE SET
		    note   = COALESCE(NULLIF(EXCLUDED.note, ''),   beta_interest.note),
		    source = COALESCE(NULLIF(EXCLUDED.source, ''), beta_interest.source)`,
		email, note, source)
	if err != nil {
		s.Logger.Error("beta_interest insert", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	s.Logger.Info("beta_interest received", "email", email, "source", source)
	writeJSON(w, http.StatusCreated, map[string]any{"status": "received"})
}

type betaInterestDTO struct {
	Email     string     `json:"email"`
	Note      string     `json:"note,omitempty"`
	Source    string     `json:"source,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	InvitedAt *time.Time `json:"invited_at,omitempty"`
}

// handleBetaInterestList is admin-only — shows pending first, then invited.
func (s *Server) handleBetaInterestList(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Pool.Query(r.Context(), `
		SELECT email, note, source, created_at, invited_at
		FROM beta_interest
		ORDER BY invited_at NULLS FIRST, created_at DESC`)
	if err != nil {
		s.Logger.Error("beta_interest list", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()
	out := []betaInterestDTO{}
	for rows.Next() {
		var d betaInterestDTO
		var note, source *string
		if err := rows.Scan(&d.Email, &note, &source, &d.CreatedAt, &d.InvitedAt); err != nil {
			s.Logger.Error("beta_interest scan", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		if note != nil {
			d.Note = *note
		}
		if source != nil {
			d.Source = *source
		}
		out = append(out, d)
	}
	writeJSON(w, http.StatusOK, out)
}

// handleBetaInterestPromote moves an interest row into invited_emails so the
// person can sign in, and stamps invited_at on the interest row so it stops
// showing up in the pending list. Wrapped in a tx so we never end up with
// the interest row promoted but the invite row missing.
func (s *Server) handleBetaInterestPromote(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	email := strings.ToLower(strings.TrimSpace(r.PathValue("email")))
	if email == "" {
		writeJSONError(w, http.StatusBadRequest, "missing email")
		return
	}

	tx, err := s.Pool.Begin(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer tx.Rollback(r.Context())

	var note *string
	err = tx.QueryRow(r.Context(),
		`SELECT note FROM beta_interest WHERE email = $1 AND invited_at IS NULL`,
		email).Scan(&note)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSONError(w, http.StatusNotFound, "no pending request for that email")
		return
	}
	if err != nil {
		s.Logger.Error("beta_interest promote: lookup", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}

	inviteNote := ""
	if note != nil {
		inviteNote = *note
	}
	if err := s.Auth.AddInvite(r.Context(), email, inviteNote, u.ID); err != nil {
		s.Logger.Info("beta_interest promote: add invite rejected", "err", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := tx.Exec(r.Context(),
		`UPDATE beta_interest SET invited_at = now() WHERE email = $1`, email,
	); err != nil {
		s.Logger.Error("beta_interest promote: stamp", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	s.Logger.Info("beta_interest promoted", "email", email, "by_user_id", u.ID)
	writeJSON(w, http.StatusOK, map[string]any{"status": "invited"})
}
