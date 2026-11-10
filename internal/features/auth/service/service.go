package auth_service

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
	auth_oauth_google "todo-api/internal/features/auth/oauth/google"
)

type AuthService struct {
	provider     OAuthGoogle
	identityRepo IdentityRepo
	userRepo     UserRepo
}

type OAuthGoogle interface {
	AuthCodeURL(state string) string
	CompleteLogin(ctx context.Context, code string) (*auth_oauth_google.ExternalUser, error)
}

type IdentityRepo interface {
	GetUserID(
		ctx context.Context,
		provider, subject string,
	) (int, error)

	CreateUserWithIdentity(
		ctx context.Context,
		user domain.User,
		provider, subject, email string,
	) (domain.User, error)
}

type UserRepo interface {
	GetUser(ctx context.Context, id int) (domain.User, error)
}

func NewAuthService(
	googleProvider OAuthGoogle,
	identityRepo IdentityRepo,
	userRepo UserRepo,
) *AuthService {
	return &AuthService{
		provider:     googleProvider,
		identityRepo: identityRepo,
		userRepo:     userRepo,
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

	userID, err := s.identityRepo.GetUserID(ctx, provider, raw.ID)
	if err != nil {
		return domain.User{}, fmt.Errorf("not found: %w", err)
	}

	var userDomain domain.User

	if userID == 0 {
		userDomain, err = s.identityRepo.CreateUserWithIdentity(
			ctx,
			domain.User{
				ID:          domain.UninitializedID,
				Version:     domain.UninitializedID,
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
	} else {
		userDomain, err = s.userRepo.GetUser(ctx, userID)
		if err != nil {
			return domain.User{}, fmt.Errorf("failed to get user: %w", err)
		}
	}

	return userDomain, nil
}
