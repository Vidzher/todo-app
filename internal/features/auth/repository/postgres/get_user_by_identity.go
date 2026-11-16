package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
	core_postgres_pool "todo-api/internal/core/repository/postgres/pool"
	users_postgres_repository "todo-api/internal/features/users/repository/postgres"
)

func (r *AuthRepository) GetUserByIdentity(
	ctx context.Context,
	provider, subject string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		select u.id, u.version, u.full_name, u.phone_number
		from todoapp.users u
		join todoapp.user_identities ui on u.id = ui.user_id
		where ui.provider = $1 and ui.subject=$2;
	`

	row := r.pool.QueryRow(ctx, query, provider, subject)

	var userModel users_postgres_repository.UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if errors.Is(err, core_postgres_pool.ErrNoRows) {
		return domain.User{}, fmt.Errorf(
			"not found with provider=%s and subject=%s: %w",
			provider,
			subject,
			core_errors.ErrNotFound,
		)
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("scan identity: %w", err)
	}

	return domain.User{
		ID:          userModel.ID,
		Version:     userModel.Version,
		FullName:    userModel.FullName,
		PhoneNumber: userModel.PhoneNumber,
	}, nil
}
