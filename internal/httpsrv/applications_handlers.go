package httpsrv

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type application struct {
	ID         int64      `json:"id"`
	Company    string     `json:"company"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	Source     *string    `json:"source"`
	JDURL      *string    `json:"jd_url"`
	Location   *string    `json:"location"`
	SalaryNote *string    `json:"salary_note"`
	CVVariant  *string    `json:"cv_variant"`
	Notes      *string    `json:"notes"`
	AppliedAt  *time.Time `json:"applied_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

var validStatus = map[string]bool{
	"wishlist": true, "applied": true, "screen": true,
	"interview": true, "offer": true, "rejected": true, "withdrawn": true,
}

func (s *Server) handleApplicationsList(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	rows, err := s.Pool.Query(r.Context(), `
		SELECT id, company, role, status, source, jd_url, location,
		       salary_note, cv_variant, notes, applied_at, created_at, updated_at
		FROM applications
		WHERE user_id = $1
		ORDER BY COALESCE(applied_at, created_at) DESC`, u.ID)
	if err != nil {
		s.Logger.Error("list applications", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()

	out := []application{}
	for rows.Next() {
		var a application
		if err := rows.Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.Source,
			&a.JDURL, &a.Location, &a.SalaryNote, &a.CVVariant, &a.Notes,
			&a.AppliedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			s.Logger.Error("scan application", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		out = append(out, a)
	}
	writeJSON(w, http.StatusOK, out)
}

type applicationInput struct {
	Company    string     `json:"company"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	Source     *string    `json:"source"`
	JDURL      *string    `json:"jd_url"`
	Location   *string    `json:"location"`
	SalaryNote *string    `json:"salary_note"`
	CVVariant  *string    `json:"cv_variant"`
	Notes      *string    `json:"notes"`
	AppliedAt  *time.Time `json:"applied_at"`
}

func (s *Server) handleApplicationCreate(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	var in applicationInput
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	in.Company = strings.TrimSpace(in.Company)
	in.Role = strings.TrimSpace(in.Role)
	if in.Company == "" || in.Role == "" {
		writeJSONError(w, http.StatusBadRequest, "company and role are required")
		return
	}
	if in.Status == "" {
		in.Status = "applied"
	}
	if !validStatus[in.Status] {
		writeJSONError(w, http.StatusBadRequest, "invalid status")
		return
	}
	if in.AppliedAt == nil && in.Status != "wishlist" {
		now := time.Now()
		in.AppliedAt = &now
	}

	var a application
	err := s.Pool.QueryRow(r.Context(), `
		INSERT INTO applications (user_id, company, role, status, source, jd_url,
		    location, salary_note, cv_variant, notes, applied_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, company, role, status, source, jd_url, location,
		    salary_note, cv_variant, notes, applied_at, created_at, updated_at`,
		u.ID, in.Company, in.Role, in.Status, in.Source, in.JDURL,
		in.Location, in.SalaryNote, in.CVVariant, in.Notes, in.AppliedAt,
	).Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.Source, &a.JDURL,
		&a.Location, &a.SalaryNote, &a.CVVariant, &a.Notes,
		&a.AppliedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		s.Logger.Error("insert application", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

func (s *Server) handleApplicationGet(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	var a application
	err = s.Pool.QueryRow(r.Context(), `
		SELECT id, company, role, status, source, jd_url, location,
		    salary_note, cv_variant, notes, applied_at, created_at, updated_at
		FROM applications WHERE id = $1 AND user_id = $2`, id, u.ID,
	).Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.Source, &a.JDURL,
		&a.Location, &a.SalaryNote, &a.CVVariant, &a.Notes,
		&a.AppliedAt, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		s.Logger.Error("get application", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (s *Server) handleApplicationUpdate(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	var in applicationInput
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if in.Status != "" && !validStatus[in.Status] {
		writeJSONError(w, http.StatusBadRequest, "invalid status")
		return
	}

	// Partial update via COALESCE — only non-empty fields overwrite.
	var a application
	err = s.Pool.QueryRow(r.Context(), `
		UPDATE applications SET
		    company    = COALESCE(NULLIF($3, ''), company),
		    role       = COALESCE(NULLIF($4, ''), role),
		    status     = COALESCE(NULLIF($5, ''), status),
		    source     = COALESCE($6, source),
		    jd_url     = COALESCE($7, jd_url),
		    location   = COALESCE($8, location),
		    salary_note= COALESCE($9, salary_note),
		    cv_variant = COALESCE($10, cv_variant),
		    notes      = COALESCE($11, notes),
		    applied_at = COALESCE($12, applied_at),
		    updated_at = now()
		WHERE id = $1 AND user_id = $2
		RETURNING id, company, role, status, source, jd_url, location,
		    salary_note, cv_variant, notes, applied_at, created_at, updated_at`,
		id, u.ID, in.Company, in.Role, in.Status, in.Source, in.JDURL,
		in.Location, in.SalaryNote, in.CVVariant, in.Notes, in.AppliedAt,
	).Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.Source, &a.JDURL,
		&a.Location, &a.SalaryNote, &a.CVVariant, &a.Notes,
		&a.AppliedAt, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		s.Logger.Error("update application", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (s *Server) handleApplicationDelete(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	tag, err := s.Pool.Exec(r.Context(), `DELETE FROM applications WHERE id = $1 AND user_id = $2`, id, u.ID)
	if err != nil {
		s.Logger.Error("delete application", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if tag.RowsAffected() == 0 {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
