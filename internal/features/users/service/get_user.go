package users_service

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	if id <= 0 {
		return domain.User{}, fmt.Errorf(
			"id must be greater than 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
