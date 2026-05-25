package httpsrv

import (
	"context"
	"net/http"
	"time"
)

// demoApp is one row of the seed fixture. AppliedDaysAgo is converted to a
// real timestamp at insert so the dashboard's date-derived copy looks alive.
type demoApp struct {
	Company        string
	Role           string
	Status         string
	Source         string
	Location       string
	SalaryNote     string
	CVVariant      string
	JDURL          string
	Notes          string // gets "[demo] " prefix automatically
	AppliedDaysAgo int    // 0 = today; negative is in the future (don't use)
}

// demoSeed is the fixture set. Designed so the dashboard's Today / This week
// sections, the Briefer, the Board, and the Funnel all have something to
// render. Mix of statuses, sources, dates spread over ~55 days.
var demoSeed = []demoApp{
	// One live offer — drives the Briefer's primary path.
	{Company: "Anthropic", Role: "Member of Technical Staff, Applied AI", Status: "offer",
		Source: "Referral", Location: "San Francisco / Remote",
		SalaryNote: "$340k base + equity", CVVariant: "v3-ai-focus",
		JDURL: "https://www.anthropic.com/jobs", Notes: "Loop went well. Verbal offer Tue.",
		AppliedDaysAgo: 38},

	// Two live interviews — populate Today + This week.
	{Company: "Stripe", Role: "Staff Engineer, Platform", Status: "interview",
		Source: "Referral", Location: "Remote (US)",
		SalaryNote: "$280-320k base", CVVariant: "v2-infra",
		Notes: "Onsite Thursday. 4 panels — sys design, IC code, behav, bar raiser.",
		AppliedDaysAgo: 24},
	{Company: "Linear", Role: "Senior Product Engineer", Status: "interview",
		Source: "LinkedIn", Location: "Remote (EU)",
		SalaryNote: "$210k base", CVVariant: "v3-ai-focus",
		Notes: "Pair-programming round Mon with Karri.",
		AppliedDaysAgo: 17},

	// Three screens — middle of the funnel.
	{Company: "Vercel", Role: "Senior Software Engineer, Edge", Status: "screen",
		Source: "Greenhouse", Location: "Remote",
		CVVariant: "v2-infra",
		Notes: "Recruiter screen done. Tech screen scheduled next week.",
		AppliedDaysAgo: 11},
	{Company: "Figma", Role: "Senior Engineer, Multiplayer", Status: "screen",
		Source: "Referral", Location: "San Francisco",
		CVVariant: "v3-ai-focus",
		Notes: "Recruiter call went well — felt very prepared.",
		AppliedDaysAgo: 9},
	{Company: "Notion", Role: "Staff Engineer, AI", Status: "screen",
		Source: "Cold email", Location: "San Francisco",
		Notes: "Cold-emailed VPE, got intro to recruiter.",
		AppliedDaysAgo: 6},

	// Five applied — recent activity at the top of the table.
	{Company: "Cursor", Role: "Founding Engineer, Agents", Status: "applied",
		Source: "X DM", Location: "San Francisco",
		Notes: "Sent DM with my agent-eval blog post.",
		AppliedDaysAgo: 4},
	{Company: "Mistral AI", Role: "Senior Research Engineer", Status: "applied",
		Source: "Company site", Location: "Paris / Remote",
		CVVariant: "v3-ai-focus",
		Notes: "Heard back from recruiter same day — promising.",
		AppliedDaysAgo: 3},
	{Company: "Perplexity", Role: "Staff Engineer, Search", Status: "applied",
		Source: "LinkedIn", Location: "San Francisco",
		AppliedDaysAgo: 2},
	{Company: "Granola", Role: "Founding Backend Engineer", Status: "applied",
		Source: "Referral", Location: "London",
		Notes: "Intro via Tom.",
		AppliedDaysAgo: 2},
	{Company: "OpenAI", Role: "Member of Technical Staff", Status: "applied",
		Source: "Company site", Location: "San Francisco",
		CVVariant: "v3-ai-focus",
		AppliedDaysAgo: 1},

	// Two wishlist — pipeline future.
	{Company: "Arc (The Browser Company)", Role: "Senior Engineer, AI Browser", Status: "wishlist",
		Source: "LinkedIn", Location: "Remote",
		Notes: "Watching for the right AI-focused opening.",
		AppliedDaysAgo: 0},
	{Company: "Replit", Role: "Staff Engineer, Agents", Status: "wishlist",
		Source: "Company site", Location: "Remote",
		AppliedDaysAgo: 0},

	// Closed loops — the bottom of the table, useful for funnel realism.
	{Company: "Plaid", Role: "Senior Engineer, Identity", Status: "rejected",
		Source: "LinkedIn", Location: "Remote (US)",
		CVVariant: "v1-generalist",
		Notes: "Failed at sys design. Need to drill caching patterns.",
		AppliedDaysAgo: 52},
	{Company: "Fly.io", Role: "Senior Engineer, Platform", Status: "withdrawn",
		Source: "Cold email", Location: "Remote",
		Notes: "Pulled after the Anthropic loop went well.",
		AppliedDaysAgo: 47},
}

const demoNotePrefix = "[demo] "

// handleDemoSeed inserts the demoSeed fixture for the calling user. Refuses
// to run if any demo rows already exist for this user — call DELETE first to
// reseed. Admin-only.
func (s *Server) handleDemoSeed(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())

	var existing int
	if err := s.Pool.QueryRow(r.Context(),
		`SELECT count(*) FROM applications WHERE user_id = $1 AND notes LIKE $2`,
		u.ID, demoNotePrefix+"%",
	).Scan(&existing); err != nil {
		s.Logger.Error("demo seed: count existing", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	if existing > 0 {
		writeJSONError(w, http.StatusConflict,
			"demo data already present — clear it first if you want to reseed")
		return
	}

	tx, err := s.Pool.Begin(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer tx.Rollback(context.Background())

	now := time.Now()
	for _, d := range demoSeed {
		appliedAt := now.AddDate(0, 0, -d.AppliedDaysAgo)
		var appliedArg any = appliedAt
		if d.Status == "wishlist" {
			appliedArg = nil
		}
		notes := demoNotePrefix + d.Notes
		var (
			source     any = nullIfEmpty(d.Source)
			location   any = nullIfEmpty(d.Location)
			salaryNote any = nullIfEmpty(d.SalaryNote)
			cvVariant  any = nullIfEmpty(d.CVVariant)
			jdURL      any = nullIfEmpty(d.JDURL)
		)
		if _, err := tx.Exec(r.Context(), `
			INSERT INTO applications (user_id, company, role, status, source, jd_url,
			    location, salary_note, cv_variant, notes, applied_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			u.ID, d.Company, d.Role, d.Status, source, jdURL,
			location, salaryNote, cvVariant, notes, appliedArg,
		); err != nil {
			s.Logger.Error("demo seed: insert", "err", err, "company", d.Company)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
	}
	if err := tx.Commit(r.Context()); err != nil {
		s.Logger.Error("demo seed: commit", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	s.Logger.Info("demo seed inserted", "user_id", u.ID, "rows", len(demoSeed))
	writeJSON(w, http.StatusCreated, map[string]any{"inserted": len(demoSeed)})
}

// handleDemoClear deletes every application whose notes start with the demo
// prefix for the calling user. Cascades to dossiers via FK if any were
// generated against demo rows. Admin-only.
func (s *Server) handleDemoClear(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	tag, err := s.Pool.Exec(r.Context(),
		`DELETE FROM applications WHERE user_id = $1 AND notes LIKE $2`,
		u.ID, demoNotePrefix+"%",
	)
	if err != nil {
		s.Logger.Error("demo clear", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	s.Logger.Info("demo clear", "user_id", u.ID, "rows", tag.RowsAffected())
	writeJSON(w, http.StatusOK, map[string]any{"deleted": tag.RowsAffected()})
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
