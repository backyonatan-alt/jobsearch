package httpsrv

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type application struct {
	ID                    int64           `json:"id"`
	Company               string          `json:"company"`
	Role                  string          `json:"role"`
	Status                string          `json:"status"`
	Source                *string         `json:"source"`
	JDURL                 *string         `json:"jd_url"`
	JDText                *string         `json:"jd_text"`
	Location              *string         `json:"location"`
	SalaryNote            *string         `json:"salary_note"`
	CVVariant             *string         `json:"cv_variant"`
	Notes                 *string         `json:"notes"`
	HiringManagerName     *string         `json:"hiring_manager_name"`
	HiringManagerLinkedIn *string         `json:"hiring_manager_linkedin"`
	RecruiterName         *string         `json:"recruiter_name"`
	RecruiterEmail        *string         `json:"recruiter_email"`
	RecruiterLinkedIn     *string         `json:"recruiter_linkedin"`
	Pipeline              json.RawMessage `json:"pipeline,omitempty"`
	AppliedAt             *time.Time      `json:"applied_at"`
	LastFollowUpAt        *time.Time      `json:"last_follow_up_at"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

var validStatus = map[string]bool{
	"wishlist": true, "applied": true, "screen": true,
	"interview": true, "offer": true, "rejected": true, "withdrawn": true,
	// "closed" = the company cancelled the req mid-process (neutral terminal —
	// not a rejection, not a candidate withdrawal).
	"closed": true,
}

func (s *Server) handleApplicationsList(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	rows, err := s.Pool.Query(r.Context(), `
		SELECT a.id, a.company, a.role, a.status, a.source, a.jd_url, a.jd_text, a.location,
		       a.salary_note, a.cv_variant, a.notes, a.hiring_manager_name,
		       a.hiring_manager_linkedin, a.recruiter_name, a.recruiter_email,
		       a.recruiter_linkedin, a.applied_at,
		       (SELECT MAX(f.occurred_at) FROM follow_ups f WHERE f.application_id = a.id) AS last_follow_up_at,
		       a.created_at, a.updated_at
		FROM applications a
		WHERE a.user_id = $1
		ORDER BY COALESCE(a.applied_at, a.created_at) DESC`, u.ID)
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
			&a.JDURL, &a.JDText, &a.Location, &a.SalaryNote, &a.CVVariant, &a.Notes,
			&a.HiringManagerName, &a.HiringManagerLinkedIn,
			&a.RecruiterName, &a.RecruiterEmail, &a.RecruiterLinkedIn,
			&a.AppliedAt, &a.LastFollowUpAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			s.Logger.Error("scan application", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		out = append(out, a)
	}
	writeJSON(w, http.StatusOK, out)
}

type applicationInput struct {
	Company               string     `json:"company"`
	Role                  string     `json:"role"`
	Status                string     `json:"status"`
	Source                *string    `json:"source"`
	JDURL                 *string    `json:"jd_url"`
	JDText                *string    `json:"jd_text"`
	Location              *string    `json:"location"`
	SalaryNote            *string    `json:"salary_note"`
	CVVariant             *string    `json:"cv_variant"`
	Notes                 *string    `json:"notes"`
	HiringManagerName     *string    `json:"hiring_manager_name"`
	HiringManagerLinkedIn *string    `json:"hiring_manager_linkedin"`
	RecruiterName         *string    `json:"recruiter_name"`
	RecruiterEmail        *string    `json:"recruiter_email"`
	RecruiterLinkedIn     *string    `json:"recruiter_linkedin"`
	AppliedAt             *time.Time `json:"applied_at"`
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
		    jd_text, location, salary_note, cv_variant, notes, hiring_manager_name,
		    hiring_manager_linkedin, recruiter_name, recruiter_email,
		    recruiter_linkedin, applied_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		RETURNING id, company, role, status, source, jd_url, jd_text, location,
		    salary_note, cv_variant, notes, hiring_manager_name,
		    hiring_manager_linkedin, recruiter_name, recruiter_email,
		    recruiter_linkedin, applied_at, created_at, updated_at`,
		u.ID, in.Company, in.Role, in.Status, in.Source, in.JDURL,
		in.JDText, in.Location, in.SalaryNote, in.CVVariant, in.Notes,
		in.HiringManagerName, in.HiringManagerLinkedIn,
		in.RecruiterName, in.RecruiterEmail, in.RecruiterLinkedIn, in.AppliedAt,
	).Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.Source, &a.JDURL, &a.JDText,
		&a.Location, &a.SalaryNote, &a.CVVariant, &a.Notes,
		&a.HiringManagerName, &a.HiringManagerLinkedIn,
		&a.RecruiterName, &a.RecruiterEmail, &a.RecruiterLinkedIn,
		&a.AppliedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		s.Logger.Error("insert application", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusCreated, a)

	// Activation milestone: first *real* application (demo-seeded rows excluded).
	// Emitted after the response so it never adds latency to the create.
	var realCount int
	if err := s.Pool.QueryRow(r.Context(),
		`SELECT count(*) FROM applications WHERE user_id = $1 AND (notes IS NULL OR notes NOT LIKE $2)`,
		u.ID, demoNotePrefix+"%",
	).Scan(&realCount); err == nil && realCount == 1 {
		s.logEvent(r.Context(), u.ID, "first_application", nil)
	}
}

type importRow struct {
	Company    string     `json:"company"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	Source     string     `json:"source"`
	Location   string     `json:"location"`
	SalaryNote string     `json:"salary_note"`
	Notes      string     `json:"notes"`
	AppliedAt  *time.Time `json:"applied_at"`
}

type importRequest struct {
	Applications []importRow `json:"applications"`
}

const maxImportRows = 500

// POST /api/applications/import — bulk-create applications from a pasted
// spreadsheet. Rows missing both company and role are skipped; unknown statuses
// fall back to "applied". Returns how many were created and how many skipped.
func (s *Server) handleApplicationsImport(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	var in importRequest
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(in.Applications) == 0 {
		writeJSONError(w, http.StatusBadRequest, "no rows to import")
		return
	}
	if len(in.Applications) > maxImportRows {
		writeJSONError(w, http.StatusRequestEntityTooLarge, "too many rows (500 max per import)")
		return
	}

	created, skipped := 0, 0
	for _, row := range in.Applications {
		company := strings.TrimSpace(row.Company)
		role := strings.TrimSpace(row.Role)
		if company == "" || role == "" {
			skipped++
			continue
		}
		status := strings.ToLower(strings.TrimSpace(row.Status))
		if !validStatus[status] {
			status = "applied"
		}
		if _, err := s.Pool.Exec(r.Context(), `
			INSERT INTO applications (user_id, company, role, status, source, location,
			    salary_note, notes, applied_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			u.ID, company, role, status,
			nilIfEmpty(row.Source), nilIfEmpty(row.Location),
			nilIfEmpty(row.SalaryNote), nilIfEmpty(row.Notes), row.AppliedAt,
		); err != nil {
			s.Logger.Error("import insert", "err", err)
			skipped++
			continue
		}
		created++
	}

	s.logEvent(r.Context(), u.ID, "bulk_import", map[string]any{"created": created, "skipped": skipped})
	writeJSON(w, http.StatusOK, map[string]int{"created": created, "skipped": skipped})
}

func nilIfEmpty(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
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
		SELECT a.id, a.company, a.role, a.status, a.source, a.jd_url, a.jd_text, a.location,
		    a.salary_note, a.cv_variant, a.notes, a.hiring_manager_name,
		    a.hiring_manager_linkedin, a.recruiter_name, a.recruiter_email,
		    a.recruiter_linkedin, a.pipeline, a.applied_at,
		    (SELECT MAX(f.occurred_at) FROM follow_ups f WHERE f.application_id = a.id) AS last_follow_up_at,
		    a.created_at, a.updated_at
		FROM applications a WHERE a.id = $1 AND a.user_id = $2`, id, u.ID,
	).Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.Source, &a.JDURL, &a.JDText,
		&a.Location, &a.SalaryNote, &a.CVVariant, &a.Notes,
		&a.HiringManagerName, &a.HiringManagerLinkedIn,
		&a.RecruiterName, &a.RecruiterEmail, &a.RecruiterLinkedIn, &a.Pipeline,
		&a.AppliedAt, &a.LastFollowUpAt, &a.CreatedAt, &a.UpdatedAt)
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
		    company                 = COALESCE(NULLIF($3, ''), company),
		    role                    = COALESCE(NULLIF($4, ''), role),
		    status                  = COALESCE(NULLIF($5, ''), status),
		    source                  = COALESCE($6, source),
		    jd_url                  = COALESCE($7, jd_url),
		    location                = COALESCE($8, location),
		    salary_note             = COALESCE($9, salary_note),
		    cv_variant              = COALESCE($10, cv_variant),
		    notes                   = COALESCE($11, notes),
		    applied_at              = COALESCE($12, applied_at),
		    hiring_manager_name     = COALESCE($13, hiring_manager_name),
		    hiring_manager_linkedin = COALESCE($14, hiring_manager_linkedin),
		    jd_text                 = COALESCE($15, jd_text),
		    recruiter_name          = COALESCE($16, recruiter_name),
		    recruiter_email         = COALESCE($17, recruiter_email),
		    recruiter_linkedin      = COALESCE($18, recruiter_linkedin),
		    updated_at              = now()
		WHERE id = $1 AND user_id = $2
		RETURNING id, company, role, status, source, jd_url, jd_text, location,
		    salary_note, cv_variant, notes, hiring_manager_name,
		    hiring_manager_linkedin, recruiter_name, recruiter_email,
		    recruiter_linkedin, applied_at, created_at, updated_at`,
		id, u.ID, in.Company, in.Role, in.Status, in.Source, in.JDURL,
		in.Location, in.SalaryNote, in.CVVariant, in.Notes, in.AppliedAt,
		in.HiringManagerName, in.HiringManagerLinkedIn,
		in.JDText, in.RecruiterName, in.RecruiterEmail, in.RecruiterLinkedIn,
	).Scan(&a.ID, &a.Company, &a.Role, &a.Status, &a.Source, &a.JDURL, &a.JDText,
		&a.Location, &a.SalaryNote, &a.CVVariant, &a.Notes,
		&a.HiringManagerName, &a.HiringManagerLinkedIn,
		&a.RecruiterName, &a.RecruiterEmail, &a.RecruiterLinkedIn,
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

type pipelineStage struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type pipelineInput struct {
	Stages []pipelineStage `json:"stages"`
}

const (
	maxPipelineStages = 30
	maxStageNameRunes = 80
)

// PUT /api/applications/{id}/pipeline — replace the per-app stage list. The
// whole array is sent each time (the client owns add/rename/reorder/toggle),
// so we just sanitize and store it. Returns the saved stages.
func (s *Server) handlePipelineUpdate(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad id")
		return
	}
	var in pipelineInput
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	stages := make([]pipelineStage, 0, len(in.Stages))
	for _, st := range in.Stages {
		name := strings.TrimSpace(st.Name)
		if name == "" {
			continue
		}
		if len([]rune(name)) > maxStageNameRunes {
			name = string([]rune(name)[:maxStageNameRunes])
		}
		stages = append(stages, pipelineStage{Name: name, Done: st.Done})
		if len(stages) >= maxPipelineStages {
			break
		}
	}
	raw, err := json.Marshal(stages)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}

	tag, err := s.Pool.Exec(r.Context(),
		`UPDATE applications SET pipeline = $3, updated_at = now() WHERE id = $1 AND user_id = $2`,
		id, u.ID, raw)
	if err != nil {
		s.Logger.Error("pipeline update", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if tag.RowsAffected() == 0 {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"pipeline": stages})
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
