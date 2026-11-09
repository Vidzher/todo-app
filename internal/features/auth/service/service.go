package auth_service

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
	auth_oauth_google "todo-api/internal/features/auth/oauth/google"
)

type AuthService struct {
	provider OAuthGoogle
}

type OAuthGoogle interface {
	AuthCodeURL(state string) string
	CompleteLogin(ctx context.Context, code string) (*auth_oauth_google.ExternalUser, error)
}

func NewAuthService(googleProvider OAuthGoogle) *AuthService {
	return &AuthService{provider: googleProvider}
}

func (s *AuthService) AuthCodeURL(state string) string {
	return s.provider.AuthCodeURL(state)
}

// TODO: fix
func (s *AuthService) CompleteLogin(ctx context.Context, code string) (any, error) {
	raw, err := s.provider.CompleteLogin(ctx, code)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to login via google: %w", err)
	}

	return raw, nil
}
