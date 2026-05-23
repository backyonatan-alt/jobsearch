package httpsrv

import (
	"context"
	"net/http"
	"time"

	"github.com/backyonatan-alt/jobsearch/internal/auth"
)

type ctxKey int

const userCtxKey ctxKey = 1

func (s *Server) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, r)
		s.Logger.Info("http",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"dur_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (s *Server) requireUser(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(s.Cfg.SessionCookieName)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "not signed in")
			return
		}
		u, ok, err := s.Auth.UserBySession(r.Context(), c.Value)
		if err != nil {
			s.Logger.Error("session lookup", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "session expired")
			return
		}
		ctx := context.WithValue(r.Context(), userCtxKey, u)
		h(w, r.WithContext(ctx))
	}
}

func userFromCtx(ctx context.Context) (auth.User, bool) {
	u, ok := ctx.Value(userCtxKey).(auth.User)
	return u, ok
}

func (s *Server) requireAdmin(h http.HandlerFunc) http.HandlerFunc {
	return s.requireUser(func(w http.ResponseWriter, r *http.Request) {
		u, _ := userFromCtx(r.Context())
		if !u.IsAdmin {
			writeJSONError(w, http.StatusForbidden, "admin only")
			return
		}
		h(w, r)
	})
}

func (s *Server) serveStaticFile(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := s.Static.Open("/" + name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			http.Error(w, "stat", http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, name, fi.ModTime(), f)
	}
}
