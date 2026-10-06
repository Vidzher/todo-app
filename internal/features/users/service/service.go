package users_service

import (
	"context"
	"todo-api/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

func NewUsersService(repo UsersRepository) *UsersService {
	return &UsersService{usersRepository: repo}
}
