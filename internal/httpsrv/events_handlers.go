package httpsrv

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// First-party analytics event log. Client UI events arrive via POST /api/events;
// server-side AI-moment events are emitted directly from the handlers. PII rule:
// never put email / company / role / JD or dossier text in props — only enums,
// booleans, counts, durations. Enforced by convention at every call site.

var eventNameRe = regexp.MustCompile(`^[a-z][a-z0-9_]{0,39}$`)

const maxEventPropsBytes = 4 * 1024

// logEvent inserts an analytics event. Best-effort: a failure is logged but
// never surfaced to the caller (analytics must not break a user action).
func (s *Server) logEvent(ctx context.Context, userID int64, name string, props map[string]any) {
	var raw []byte
	if len(props) > 0 {
		b, err := json.Marshal(props)
		if err != nil {
			s.Logger.Error("event marshal", "err", err, "event", name)
			return
		}
		raw = b
	}
	s.logEventRaw(ctx, userID, name, raw)
}

func (s *Server) logEventRaw(ctx context.Context, userID int64, name string, props []byte) {
	if len(props) == 0 {
		props = []byte("{}")
	}
	// Detach from the request context so a client disconnect after the response
	// doesn't cancel the insert.
	bg, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
	defer cancel()
	if _, err := s.Pool.Exec(bg,
		`INSERT INTO events (user_id, name, props) VALUES ($1, $2, $3)`,
		userID, name, props,
	); err != nil {
		s.Logger.Error("event insert", "err", err, "event", name)
	}
}

type eventInput struct {
	Name  string          `json:"name"`
	Props json.RawMessage `json:"props"`
}

// handleEventCreate ingests a single client UI event. Authed; the user is taken
// from the session, never the body.
func (s *Server) handleEventCreate(w http.ResponseWriter, r *http.Request) {
	u, _ := userFromCtx(r.Context())
	var in eventInput
	if err := readJSON(r, &in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	name := strings.TrimSpace(in.Name)
	if !eventNameRe.MatchString(name) {
		writeJSONError(w, http.StatusBadRequest, "invalid event name")
		return
	}
	if len(in.Props) > maxEventPropsBytes {
		writeJSONError(w, http.StatusRequestEntityTooLarge, "props too large")
		return
	}
	var props []byte
	if len(in.Props) > 0 {
		// Validate it's well-formed JSON before storing.
		if !json.Valid(in.Props) {
			writeJSONError(w, http.StatusBadRequest, "props must be valid JSON")
			return
		}
		props = in.Props
	}
	s.logEventRaw(r.Context(), u.ID, name, props)
	w.WriteHeader(http.StatusNoContent)
}

type eventDTO struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Props     json.RawMessage `json:"props"`
	CreatedAt time.Time       `json:"created_at"`
}

// handleAdminEventsList returns recent events for debugging/analysis — the
// per-user timeline that is the primary read for a small beta. Filters:
// ?user=<id|email>, ?name=<event_name>, ?limit=<n> (default 200, max 1000).
func (s *Server) handleAdminEventsList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit := 200
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > 1000 {
		limit = 1000
	}

	where := []string{}
	args := []any{}
	if user := strings.TrimSpace(q.Get("user")); user != "" {
		if id, err := strconv.ParseInt(user, 10, 64); err == nil {
			args = append(args, id)
			where = append(where, "e.user_id = $"+strconv.Itoa(len(args)))
		} else {
			args = append(args, strings.ToLower(user))
			where = append(where, "lower(u.email) = $"+strconv.Itoa(len(args)))
		}
	}
	if name := strings.TrimSpace(q.Get("name")); name != "" {
		args = append(args, name)
		where = append(where, "e.name = $"+strconv.Itoa(len(args)))
	}

	sql := `SELECT e.id, e.user_id, u.email, e.name, e.props, e.created_at
	        FROM events e JOIN users u ON u.id = e.user_id`
	if len(where) > 0 {
		sql += " WHERE " + strings.Join(where, " AND ")
	}
	args = append(args, limit)
	sql += " ORDER BY e.created_at DESC LIMIT $" + strconv.Itoa(len(args))

	rows, err := s.Pool.Query(r.Context(), sql, args...)
	if err != nil {
		s.Logger.Error("events list", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()

	out := []eventDTO{}
	for rows.Next() {
		var e eventDTO
		var props []byte
		if err := rows.Scan(&e.ID, &e.UserID, &e.Email, &e.Name, &props, &e.CreatedAt); err != nil {
			s.Logger.Error("events scan", "err", err)
			writeJSONError(w, http.StatusInternalServerError, "internal")
			return
		}
		e.Props = json.RawMessage(props)
		out = append(out, e)
	}
	writeJSON(w, http.StatusOK, out)
}
