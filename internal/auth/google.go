package auth

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

// Google holds the OAuth2 config + ID-token audience needed to drive the
// Google sign-in flow.
type Google struct {
	cfg      *oauth2.Config
	clientID string
}

func NewGoogle(clientID, clientSecret, redirectURL string) *Google {
	return &Google{
		clientID: clientID,
		cfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"openid", "email", "profile"},
		},
	}
}

// AuthCodeURL returns the URL to redirect a user to in order to start the
// Google sign-in flow. `state` should be a random per-request token also stored
// in a cookie, so the callback can verify it.
func (g *Google) AuthCodeURL(state string) string {
	return g.cfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeAndVerify takes the `code` query param from Google's callback,
// exchanges it for tokens, and returns the verified email + the profile
// picture URL from the ID token. The picture claim is optional — if Google
// doesn't include one, picture is returned as "".
func (g *Google) ExchangeAndVerify(ctx context.Context, code string) (email, picture string, err error) {
	tok, err := g.cfg.Exchange(ctx, code)
	if err != nil {
		return "", "", fmt.Errorf("exchange code: %w", err)
	}
	raw, ok := tok.Extra("id_token").(string)
	if !ok || raw == "" {
		return "", "", errors.New("no id_token in google response")
	}
	payload, err := idtoken.Validate(ctx, raw, g.clientID)
	if err != nil {
		return "", "", fmt.Errorf("validate id_token: %w", err)
	}
	e, _ := payload.Claims["email"].(string)
	verified, _ := payload.Claims["email_verified"].(bool)
	if e == "" || !verified {
		return "", "", errors.New("google account has no verified email")
	}
	pic, _ := payload.Claims["picture"].(string)
	return e, pic, nil
}
