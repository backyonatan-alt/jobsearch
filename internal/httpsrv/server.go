package httpsrv

import (
	"log/slog"
	"net/http"

	"github.com/backyonatan-alt/jobsearch/internal/auth"
	"github.com/backyonatan-alt/jobsearch/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Cfg    *config.Config
	Pool   *pgxpool.Pool
	Auth   *auth.Service
	Google *auth.Google
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
	mux.HandleFunc("GET /api/me", s.requireUser(s.handleMe))

	mux.HandleFunc("GET /api/applications", s.requireUser(s.handleApplicationsList))
	mux.HandleFunc("POST /api/applications", s.requireUser(s.handleApplicationCreate))
	mux.HandleFunc("GET /api/applications/{id}", s.requireUser(s.handleApplicationGet))
	mux.HandleFunc("PATCH /api/applications/{id}", s.requireUser(s.handleApplicationUpdate))
	mux.HandleFunc("DELETE /api/applications/{id}", s.requireUser(s.handleApplicationDelete))

	// SPA fallback: the SvelteKit static build emits a single index.html for
	// all client-routed pages, so any /app/* URL hits the same entry file.
	mux.HandleFunc("GET /app", s.serveStaticFile("index.html"))
	mux.HandleFunc("GET /app/", s.serveStaticFile("index.html"))
	mux.Handle("/", http.FileServer(s.Static))

	return s.withLogging(mux)
}
