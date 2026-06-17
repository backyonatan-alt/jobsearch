package httpsrv

import (
	"log/slog"
	"net/http"

	"github.com/backyonatan-alt/jobsearch/internal/auth"
	"github.com/backyonatan-alt/jobsearch/internal/config"
	"github.com/backyonatan-alt/jobsearch/internal/llm"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Cfg    *config.Config
	Pool   *pgxpool.Pool
	Auth   *auth.Service
	Google *auth.Google
	LLM    *llm.Client // nil when ANTHROPIC_API_KEY isn't set
	Logger *slog.Logger
	Static http.FileSystem
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /auth/google/start", s.handleGoogleStart)
	mux.HandleFunc("GET /auth/google/callback", s.handleGoogleCallback)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
	mux.HandleFunc("POST /api/beta-interest", s.handleBetaInterestCreate)
	mux.HandleFunc("GET /api/me", s.requireUser(s.handleMe))
	mux.HandleFunc("POST /api/me/onboarded", s.requireUser(s.handleMarkOnboarded))
	mux.HandleFunc("POST /api/events", s.requireUser(s.handleEventCreate))
	// Seed/clear the caller's own demo data — used by the first-run guided tour
	// to populate the screens during the tour, then clear them on finish.
	mux.HandleFunc("POST /api/me/demo-seed", s.requireUser(s.handleDemoSeed))
	mux.HandleFunc("DELETE /api/me/demo-seed", s.requireUser(s.handleDemoClear))
	mux.HandleFunc("GET /api/applications", s.requireUser(s.handleApplicationsList))
	mux.HandleFunc("POST /api/applications", s.requireUser(s.handleApplicationCreate))
	mux.HandleFunc("POST /api/applications/parse", s.requireUser(s.handleApplicationParse))
	mux.HandleFunc("POST /api/applications/import", s.requireUser(s.handleApplicationsImport))
	mux.HandleFunc("GET /api/applications/{id}", s.requireUser(s.handleApplicationGet))
	mux.HandleFunc("PATCH /api/applications/{id}", s.requireUser(s.handleApplicationUpdate))
	mux.HandleFunc("DELETE /api/applications/{id}", s.requireUser(s.handleApplicationDelete))
	mux.HandleFunc("PUT /api/applications/{id}/pipeline", s.requireUser(s.handlePipelineUpdate))
	mux.HandleFunc("GET /api/applications/{id}/dossier", s.requireUser(s.handleDossierGet))
	mux.HandleFunc("POST /api/applications/{id}/dossier/refresh", s.requireUser(s.handleDossierRefresh))

	mux.HandleFunc("GET /api/applications/{id}/interviews", s.requireUser(s.handleInterviewsList))
	mux.HandleFunc("POST /api/applications/{id}/interviews", s.requireUser(s.handleInterviewCreate))
	mux.HandleFunc("POST /api/applications/{id}/interviews/parse", s.requireUser(s.handleInterviewsParse))
	mux.HandleFunc("DELETE /api/applications/{id}/interviews/{iid}", s.requireUser(s.handleInterviewDelete))

	mux.HandleFunc("GET /api/applications/{id}/follow-ups", s.requireUser(s.handleFollowUpsList))
	mux.HandleFunc("POST /api/applications/{id}/follow-ups", s.requireUser(s.handleFollowUpCreate))
	mux.HandleFunc("DELETE /api/applications/{id}/follow-ups/{fid}", s.requireUser(s.handleFollowUpDelete))

	mux.HandleFunc("GET /api/admin/invites", s.requireAdmin(s.handleAdminInvitesList))
	mux.HandleFunc("POST /api/admin/invites", s.requireAdmin(s.handleAdminInvitesAdd))
	mux.HandleFunc("DELETE /api/admin/invites/{email}", s.requireAdmin(s.handleAdminInvitesDelete))
	mux.HandleFunc("POST /api/admin/demo-seed", s.requireAdmin(s.handleDemoSeed))
	mux.HandleFunc("DELETE /api/admin/demo-seed", s.requireAdmin(s.handleDemoClear))
	mux.HandleFunc("GET /api/admin/users", s.requireAdmin(s.handleAdminUsersList))
	mux.HandleFunc("GET /api/admin/invite-funnel", s.requireAdmin(s.handleAdminInviteFunnel))
	mux.HandleFunc("GET /api/admin/adoption", s.requireAdmin(s.handleAdminAdoption))
	mux.HandleFunc("POST /api/admin/users/{id}/prep-credits", s.requireAdmin(s.handleAdminGrantPrep))
	mux.HandleFunc("GET /api/admin/events", s.requireAdmin(s.handleAdminEventsList))
	mux.HandleFunc("GET /api/admin/beta-interest", s.requireAdmin(s.handleBetaInterestList))
	mux.HandleFunc("POST /api/admin/beta-interest/{email}/invite", s.requireAdmin(s.handleBetaInterestPromote))

	// SPA fallback: the SvelteKit static build emits a single index.html for
	// all client-routed pages, so any /app/* or /admin/* or /preview/* URL
	// hits the same entry file. The Svelte router takes over from there.
	mux.HandleFunc("GET /app", s.serveIndexHTML)
	mux.HandleFunc("GET /app/", s.serveIndexHTML)
	mux.HandleFunc("GET /admin", s.serveIndexHTML)
	mux.HandleFunc("GET /admin/", s.serveIndexHTML)
	mux.HandleFunc("GET /preview", s.serveIndexHTML)
	mux.HandleFunc("GET /preview/", s.serveIndexHTML)
	// Root serves the same shell (with GA injected); everything else under "/"
	// is a static asset (js/css/favicon) served straight off disk.
	mux.HandleFunc("GET /{$}", s.serveIndexHTML)
	mux.Handle("/", http.FileServer(s.Static))

	return s.withLogging(mux)
}
