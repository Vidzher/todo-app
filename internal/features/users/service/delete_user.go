package users_service

import (
	"context"
	"fmt"
	core_errors "todo-api/internal/core/errors"
)

func (s *UsersService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf(
			"id must be greater than 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if err := s.usersRepository.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
