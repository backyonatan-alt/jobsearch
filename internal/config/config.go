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
	MagicLinkTTL      time.Duration
	MailDriver        string
	MailFrom          string
	AnthropicAPIKey   string
	AllowedEmails     map[string]struct{} // empty map = allow any
}

func FromEnv() (*Config, error) {
	c := &Config{
		Port:              getenv("PORT", "8080"),
		BaseURL:           getenv("BASE_URL", "http://localhost:8080"),
		BrandName:         getenv("BRAND_NAME", "Pursuit"),
		DatabaseURL:       getenv("DATABASE_URL", ""),
		SessionCookieName: getenv("SESSION_COOKIE_NAME", "pursuit_session"),
		MailDriver:        getenv("MAIL_DRIVER", "log"),
		MailFrom:          getenv("MAIL_FROM", "no-reply@pursuit.local"),
		AnthropicAPIKey:   getenv("ANTHROPIC_API_KEY", ""),
	}

	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	sessionHours, err := strconv.Atoi(getenv("SESSION_TTL_HOURS", "720"))
	if err != nil {
		return nil, fmt.Errorf("SESSION_TTL_HOURS: %w", err)
	}
	c.SessionTTL = time.Duration(sessionHours) * time.Hour

	magicMin, err := strconv.Atoi(getenv("MAGIC_LINK_TTL_MINUTES", "15"))
	if err != nil {
		return nil, fmt.Errorf("MAGIC_LINK_TTL_MINUTES: %w", err)
	}
	c.MagicLinkTTL = time.Duration(magicMin) * time.Minute

	c.AllowedEmails = map[string]struct{}{}
	for _, e := range strings.Split(getenv("ALLOWED_EMAILS", ""), ",") {
		e = strings.TrimSpace(strings.ToLower(e))
		if e != "" {
			c.AllowedEmails[e] = struct{}{}
		}
	}
	return c, nil
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
