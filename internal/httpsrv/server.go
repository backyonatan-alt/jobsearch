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
	// Demo seed for the calling user — same handlers as the admin route, just
	// open to anyone signed in so a friend can populate their account with
	// realistic data in one click from the empty Today state.
	mux.HandleFunc("POST /api/me/demo-seed", s.requireUser(s.handleDemoSeed))
	mux.HandleFunc("DELETE /api/me/demo-seed", s.requireUser(s.handleDemoClear))

	mux.HandleFunc("GET /api/applications", s.requireUser(s.handleApplicationsList))
	mux.HandleFunc("POST /api/applications", s.requireUser(s.handleApplicationCreate))
	mux.HandleFunc("POST /api/applications/parse", s.requireUser(s.handleApplicationParse))
	mux.HandleFunc("GET /api/applications/{id}", s.requireUser(s.handleApplicationGet))
	mux.HandleFunc("PATCH /api/applications/{id}", s.requireUser(s.handleApplicationUpdate))
	mux.HandleFunc("DELETE /api/applications/{id}", s.requireUser(s.handleApplicationDelete))
	mux.HandleFunc("GET /api/applications/{id}/dossier", s.requireUser(s.handleDossierGet))
	mux.HandleFunc("POST /api/applications/{id}/dossier/refresh", s.requireUser(s.handleDossierRefresh))

	mux.HandleFunc("GET /api/applications/{id}/interviews", s.requireUser(s.handleInterviewsList))
	mux.HandleFunc("POST /api/applications/{id}/interviews", s.requireUser(s.handleInterviewCreate))
	mux.HandleFunc("POST /api/applications/{id}/interviews/parse", s.requireUser(s.handleInterviewsParse))
	mux.HandleFunc("DELETE /api/applications/{id}/interviews/{iid}", s.requireUser(s.handleInterviewDelete))

	mux.HandleFunc("GET /api/admin/invites", s.requireAdmin(s.handleAdminInvitesList))
	mux.HandleFunc("POST /api/admin/invites", s.requireAdmin(s.handleAdminInvitesAdd))
	mux.HandleFunc("DELETE /api/admin/invites/{email}", s.requireAdmin(s.handleAdminInvitesDelete))
	mux.HandleFunc("POST /api/admin/demo-seed", s.requireAdmin(s.handleDemoSeed))
	mux.HandleFunc("DELETE /api/admin/demo-seed", s.requireAdmin(s.handleDemoClear))
	mux.HandleFunc("GET /api/admin/beta-interest", s.requireAdmin(s.handleBetaInterestList))
	mux.HandleFunc("POST /api/admin/beta-interest/{email}/invite", s.requireAdmin(s.handleBetaInterestPromote))

	// SPA fallback: the SvelteKit static build emits a single index.html for
	// all client-routed pages, so any /app/* or /admin/* or /preview/* URL
	// hits the same entry file. The Svelte router takes over from there.
	mux.HandleFunc("GET /app", s.serveStaticFile("index.html"))
	mux.HandleFunc("GET /app/", s.serveStaticFile("index.html"))
	mux.HandleFunc("GET /admin", s.serveStaticFile("index.html"))
	mux.HandleFunc("GET /admin/", s.serveStaticFile("index.html"))
	mux.HandleFunc("GET /preview", s.serveStaticFile("index.html"))
	mux.HandleFunc("GET /preview/", s.serveStaticFile("index.html"))
	mux.Handle("/", http.FileServer(s.Static))

	return s.withLogging(mux)
}
