package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port              string
	BaseURL           string
	BrandName         string
	DatabaseURL       string
	SessionCookieName string
	SessionTTL        time.Duration

	GoogleClientID     string
	GoogleClientSecret string

	AnthropicAPIKey string
	AllowedEmails   map[string]struct{} // bootstrap fallback only; invites live in DB
	AdminEmails     map[string]struct{} // emails granted is_admin on sign-in

	GA4MeasurementID string // empty = no analytics (dev/local); set in prod env

	OpenSignup bool // OPEN_SIGNUP=true lets any Google account sign in (public beta); off = invite list
}

func FromEnv() (*Config, error) {
	c := &Config{
		Port:               getenv("PORT", "8080"),
		BaseURL:            getenv("BASE_URL", "http://localhost:8080"),
		BrandName:          getenv("BRAND_NAME", "Pursuit"),
		DatabaseURL:        getenv("DATABASE_URL", ""),
		SessionCookieName:  getenv("SESSION_COOKIE_NAME", "pursuit_session"),
		GoogleClientID:     getenv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getenv("GOOGLE_CLIENT_SECRET", ""),
		AnthropicAPIKey:    getenv("ANTHROPIC_API_KEY", ""),
		GA4MeasurementID:   strings.TrimSpace(getenv("GA4_MEASUREMENT_ID", "")),
		OpenSignup:         strings.EqualFold(strings.TrimSpace(getenv("OPEN_SIGNUP", "")), "true"),
	}

	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if c.GoogleClientID == "" || c.GoogleClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET are required")
	}

	sessionHours, err := strconv.Atoi(getenv("SESSION_TTL_HOURS", "720"))
	if err != nil {
		return nil, fmt.Errorf("SESSION_TTL_HOURS: %w", err)
	}
	c.SessionTTL = time.Duration(sessionHours) * time.Hour

	c.AllowedEmails = parseEmailSet(getenv("ALLOWED_EMAILS", ""))
	c.AdminEmails = parseEmailSet(getenv("ADMIN_EMAILS", ""))
	return c, nil
}

func parseEmailSet(s string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, e := range strings.Split(s, ",") {
		e = strings.TrimSpace(strings.ToLower(e))
		if e != "" {
			out[e] = struct{}{}
		}
	}
	return out
}

func (c *Config) IsAdminEmail(email string) bool {
	_, ok := c.AdminEmails[strings.ToLower(strings.TrimSpace(email))]
	return ok
}

func (c *Config) GoogleRedirectURL() string {
	return strings.TrimRight(c.BaseURL, "/") + "/auth/google/callback"
}

func (c *Config) EmailAllowed(email string) bool {
	if len(c.AllowedEmails) == 0 {
		return true
	}
	_, ok := c.AllowedEmails[strings.ToLower(strings.TrimSpace(email))]
	return ok
}

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
