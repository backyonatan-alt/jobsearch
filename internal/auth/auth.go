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
	ID    int64
	Email string
}

// UpsertUserAndSession creates (or updates) the user row for the given email
// and issues a new session token.
func (s *Service) UpsertUserAndSession(ctx context.Context, email, userAgent, ip string) (sessionToken string, u User, err error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") {
		return "", User{}, errors.New("invalid email")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", User{}, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO users (email) VALUES ($1)
		ON CONFLICT (email) DO UPDATE SET last_login_at = now()
		RETURNING id, email`, email).Scan(&u.ID, &u.Email)
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
		SELECT u.id, u.email
		FROM sessions s JOIN users u ON u.id = s.user_id
		WHERE s.token = $1 AND s.expires_at > now()`,
		token).Scan(&u.ID, &u.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, err
	}
	return u, true, nil
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
