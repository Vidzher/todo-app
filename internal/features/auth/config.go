package auth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Google OAuthConfig
}

type OAuthConfig struct {
	Config      *oauth2.Config
	UserInfoURL string
}

func NewConfigMust() *Config {
	return &Config{
		Google: OAuthConfig{
			Config: &oauth2.Config{
				RedirectURL:  "http://localhost:5050/api/v1/auth/google/callback",
				ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
				Scopes: []string{
					"openid",
					"profile",
					"email",
				},
				Endpoint: google.Endpoint,
			},
			UserInfoURL: "https://openidconnect.googleapis.com/v1/userinfo",
		},
	}
}
