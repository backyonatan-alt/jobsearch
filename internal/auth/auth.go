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
	pool         *pgxpool.Pool
	sessionTTL   time.Duration
	magicTTL     time.Duration
}

func NewService(pool *pgxpool.Pool, sessionTTL, magicTTL time.Duration) *Service {
	return &Service{pool: pool, sessionTTL: sessionTTL, magicTTL: magicTTL}
}

type User struct {
	ID    int64
	Email string
}

// CreateMagicLink issues a single-use token for the given email. The token is
// returned to the caller (so the mail driver can render the link); it is also
// stored hashed-by-opacity in the DB until consumed.
func (s *Service) CreateMagicLink(ctx context.Context, email string) (token string, expiresAt time.Time, err error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") {
		return "", time.Time{}, errors.New("invalid email")
	}
	token, err = randomToken(32)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt = time.Now().Add(s.magicTTL)
	_, err = s.pool.Exec(ctx, `
		INSERT INTO magic_links (token, email, expires_at) VALUES ($1, $2, $3)`,
		token, email, expiresAt)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("insert magic_link: %w", err)
	}
	return token, expiresAt, nil
}

// ConsumeMagicLink validates the token, marks it consumed, and returns a new
// session for the associated user (creating the user row if needed).
func (s *Service) ConsumeMagicLink(ctx context.Context, token, userAgent, ip string) (sessionToken string, u User, err error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", User{}, err
	}
	defer tx.Rollback(ctx)

	var email string
	var expiresAt time.Time
	var consumedAt *time.Time
	err = tx.QueryRow(ctx, `
		SELECT email, expires_at, consumed_at FROM magic_links WHERE token = $1 FOR UPDATE`,
		token).Scan(&email, &expiresAt, &consumedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", User{}, errors.New("invalid token")
	}
	if err != nil {
		return "", User{}, err
	}
	if consumedAt != nil {
		return "", User{}, errors.New("token already used")
	}
	if time.Now().After(expiresAt) {
		return "", User{}, errors.New("token expired")
	}

	if _, err := tx.Exec(ctx, `UPDATE magic_links SET consumed_at = now() WHERE token = $1`, token); err != nil {
		return "", User{}, err
	}

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

func randomToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
