package auth_oauth_google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type Provider struct {
	config      *oauth2.Config
	userInfoUrl string
	client      *http.Client
}

type googleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
}

type ExternalUser struct {
	ID    string
	Email string
	Name  string
}

func NewProvider(
	cfg *oauth2.Config,
	url string,
	c *http.Client,
) *Provider {
	return &Provider{
		config:      cfg,
		userInfoUrl: url,
		client:      c,
	}
}

func (p *Provider) AuthCodeURL(state string) string {
	return p.config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)
}

func (p *Provider) CompleteLogin(ctx context.Context, code string) (*ExternalUser, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf(
			"exchange auth code: %w",
			err,
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		p.userInfoUrl,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create user info request: %w",
			err,
		)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	res, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"execution user info request: %w",
			err,
		)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil,
			fmt.Errorf(
				"google responded with code %d",
				res.StatusCode,
			)
	}

	var user googleUser
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf(
			"decode user info: %w",
			err,
		)
	}

	return &ExternalUser{
		ID:    user.Sub,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
