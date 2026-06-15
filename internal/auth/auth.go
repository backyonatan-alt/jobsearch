package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool       *pgxpool.Pool
	sessionTTL time.Duration
}

func NewService(pool *pgxpool.Pool, sessionTTL time.Duration) *Service {
	return &Service{pool: pool, sessionTTL: sessionTTL}
}

type User struct {
	ID          int64
	Email       string
	IsAdmin     bool
	OnboardedAt *time.Time
	PictureURL  *string
}

// UpsertUserAndSession creates (or updates) the user row for the given email
// and issues a new session token. If `makeAdmin` is true the user's is_admin
// flag is set to true (used so the configured ADMIN_EMAILS automatically gain
// admin on sign-in). `pictureURL` is the Google profile picture URL from the
// ID token — refreshed on every sign-in so it stays current.
func (s *Service) UpsertUserAndSession(ctx context.Context, email, userAgent, ip, pictureURL string, makeAdmin bool) (sessionToken string, u User, err error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") {
		return "", User{}, errors.New("invalid email")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", User{}, err
	}
	defer tx.Rollback(ctx)

	// NULLIF so a missing picture claim doesn't clobber a previously stored one.
	// last_login_at is set on INSERT *and* in the conflict branch — otherwise a
	// user's very first login never records it and they show "Last seen never".
	err = tx.QueryRow(ctx, `
		INSERT INTO users (email, is_admin, picture_url, last_login_at) VALUES ($1, $2, NULLIF($3, ''), now())
		ON CONFLICT (email) DO UPDATE SET
		    last_login_at = now(),
		    is_admin = users.is_admin OR EXCLUDED.is_admin,
		    picture_url = COALESCE(NULLIF($3, ''), users.picture_url)
		RETURNING id, email, is_admin, onboarded_at, picture_url`,
		email, makeAdmin, pictureURL).Scan(&u.ID, &u.Email, &u.IsAdmin, &u.OnboardedAt, &u.PictureURL)
	if err != nil {
		return "", User{}, fmt.Errorf("upsert user: %w", err)
	}

	sessionToken, err = randomToken(32)
	if err != nil {
		return "", User{}, err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO sessions (token, user_id, expires_at, user_agent, ip)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::inet)`,
		sessionToken, u.ID, time.Now().Add(s.sessionTTL), userAgent, ip)
	if err != nil {
		return "", User{}, fmt.Errorf("create session: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", User{}, err
	}
	return sessionToken, u, nil
}

// UserBySession resolves a session token to its owning user, or returns false.
func (s *Service) UserBySession(ctx context.Context, token string) (User, bool, error) {
	if token == "" {
		return User{}, false, nil
	}
	var u User
	err := s.pool.QueryRow(ctx, `
		SELECT u.id, u.email, u.is_admin, u.onboarded_at, u.picture_url
		FROM sessions s JOIN users u ON u.id = s.user_id
		WHERE s.token = $1 AND s.expires_at > now()`,
		token).Scan(&u.ID, &u.Email, &u.IsAdmin, &u.OnboardedAt, &u.PictureURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, err
	}
	return u, true, nil
}

// EmailInvited reports whether the given email is allowed to sign in. The
// invited_emails table is the source of truth; the env-based ALLOWED_EMAILS
// list is kept as a bootstrap fallback so the first ever sign-in works on a
// fresh DB.
func (s *Service) EmailInvited(ctx context.Context, email string, envFallback map[string]struct{}) (bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return false, nil
	}
	var count int
	if err := s.pool.QueryRow(ctx,
		`SELECT count(*) FROM invited_emails WHERE lower(email) = $1`, email,
	).Scan(&count); err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	if _, ok := envFallback[email]; ok {
		return true, nil
	}
	// Empty allow-list means anyone can sign in (used in dev).
	hasDBList, err := s.invitesAny(ctx)
	if err != nil {
		return false, err
	}
	return !hasDBList && len(envFallback) == 0, nil
}

func (s *Service) invitesAny(ctx context.Context) (bool, error) {
	var any bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM invited_emails)`).Scan(&any)
	return any, err
}

// ListInvites returns every invited email plus when it was added and by whom.
func (s *Service) ListInvites(ctx context.Context) ([]Invite, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT i.email, i.invited_at, i.note, u.email
		FROM invited_emails i
		LEFT JOIN users u ON u.id = i.invited_by_user_id
		ORDER BY i.invited_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Invite{}
	for rows.Next() {
		var v Invite
		var byEmail *string
		if err := rows.Scan(&v.Email, &v.InvitedAt, &v.Note, &byEmail); err != nil {
			return nil, err
		}
		if byEmail != nil {
			v.InvitedByEmail = *byEmail
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// AddInvite inserts (or replaces, idempotently) an invited email.
func (s *Service) AddInvite(ctx context.Context, email, note string, byUserID int64) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") {
		return errors.New("invalid email")
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO invited_emails (email, invited_by_user_id, note) VALUES ($1, $2, NULLIF($3,''))
		ON CONFLICT (email) DO UPDATE SET note = EXCLUDED.note`,
		email, byUserID, note)
	return err
}

// RemoveInvite deletes an invited email. (Existing sessions are not affected;
// they remain valid until they expire or the user logs out.)
func (s *Service) RemoveInvite(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	_, err := s.pool.Exec(ctx, `DELETE FROM invited_emails WHERE lower(email) = $1`, email)
	return err
}

// SeedInvitesFromEnv copies any emails in the bootstrap env list into the DB
// the first time the table is empty. Subsequent runs are no-ops.
func (s *Service) SeedInvitesFromEnv(ctx context.Context, emails map[string]struct{}) error {
	if len(emails) == 0 {
		return nil
	}
	has, err := s.invitesAny(ctx)
	if err != nil || has {
		return err
	}
	for e := range emails {
		if _, err := s.pool.Exec(ctx,
			`INSERT INTO invited_emails (email, note) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			e, "seeded from ALLOWED_EMAILS env"); err != nil {
			return err
		}
	}
	return nil
}

// SeedAdminsFromEnv promotes any existing user matching the configured
// ADMIN_EMAILS list to admin. Idempotent — safe to call on every boot.
func (s *Service) SeedAdminsFromEnv(ctx context.Context, emails map[string]struct{}) error {
	for e := range emails {
		if _, err := s.pool.Exec(ctx,
			`UPDATE users SET is_admin = true WHERE lower(email) = $1`, e); err != nil {
			return err
		}
	}
	return nil
}

type Invite struct {
	Email          string
	InvitedAt      time.Time
	Note           *string
	InvitedByEmail string
}

// MarkOnboarded sets onboarded_at = now() for the user (idempotent — if
// already set, the row is left alone).
func (s *Service) MarkOnboarded(ctx context.Context, userID int64) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET onboarded_at = COALESCE(onboarded_at, now()) WHERE id = $1`,
		userID)
	return err
}

func (s *Service) DestroySession(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}

// RandomToken returns a URL-safe random token of nBytes bytes.
func RandomToken(nBytes int) (string, error) {
	return randomToken(nBytes)
}

func randomToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
