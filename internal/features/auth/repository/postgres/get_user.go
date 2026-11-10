package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	core_postgres_pool "todo-api/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) GetUserID(
	ctx context.Context,
	provider, subject string,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		select user_id 
		from todoapp.user_identities
		where provider=$1 and subject=$2;
	`

	row := r.pool.QueryRow(ctx, query, provider, subject)

	var userID int
	err := row.Scan(&userID)
	if errors.Is(err, core_postgres_pool.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf("scan identity: %w", err)
	}

	return userID, nil
}
