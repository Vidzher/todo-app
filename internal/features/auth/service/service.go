package auth_service

import (
	"context"
	"errors"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
	auth_oauth_google "todo-api/internal/features/auth/oauth/google"
)

type AuthService struct {
	provider     OAuthGoogle
	identityRepo IdentityRepo
}

type OAuthGoogle interface {
	AuthCodeURL(state string) string
	CompleteLogin(ctx context.Context, code string) (*auth_oauth_google.ExternalUser, error)
}

type IdentityRepo interface {
	GetUserByIdentity(
		ctx context.Context,
		provider, subject string,
	) (domain.User, error)

	CreateUserWithIdentity(
		ctx context.Context,
		user domain.User,
		provider, subject, email string,
	) (domain.User, error)
}

func NewAuthService(
	googleProvider OAuthGoogle,
	identityRepo IdentityRepo,
) *AuthService {
	return &AuthService{
		provider:     googleProvider,
		identityRepo: identityRepo,
	}
}

func (s *AuthService) AuthCodeURL(state string) string {
	return s.provider.AuthCodeURL(state)
}

// TODO: fix
// - проверить в таблице user_identities наличие записи с полученным raw.ID
// - если записи нет, создать пользователя в таблице users и user_indentities
// - если запись есть, получить запись из таблицы users по ключу user_id
// - вернуть на клиент domain.User
func (s *AuthService) CompleteLogin(ctx context.Context, code, provider string) (any, error) {
	raw, err := s.provider.CompleteLogin(ctx, code)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to login via google: %w", err)
	}

	user, err := s.identityRepo.GetUserByIdentity(ctx, provider, raw.ID)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, core_errors.ErrNotFound) {
		return domain.User{}, fmt.Errorf("find by identity: %w", err)
	}

	var userDomain domain.User

	userDomain, err = s.identityRepo.CreateUserWithIdentity(
		ctx,
		domain.User{
			ID:          domain.UninitializedID,
			Version:     domain.UninitializedVersion,
			FullName:    raw.Name,
			PhoneNumber: nil,
		},
		provider,
		raw.ID,
		raw.Email,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return userDomain, nil
}
