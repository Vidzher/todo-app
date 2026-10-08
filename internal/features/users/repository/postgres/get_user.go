package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
	core_postgres_pool "todo-api/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		select id, version, full_name, phone_number
		from todoapp.users
		where id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if errors.Is(err, core_postgres_pool.ErrNoRows) {
		return domain.User{}, fmt.Errorf(
			"user with id=%d: %v: %w",
			id,
			err,
			core_errors.ErrNotFound,
		)
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	return domain.User{
		ID:          userModel.ID,
		Version:     userModel.Version,
		FullName:    userModel.FullName,
		PhoneNumber: userModel.PhoneNumber,
	}, nil
}
