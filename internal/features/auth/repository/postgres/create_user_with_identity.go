package auth_postgres_repository

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
	users_postgres_repository "todo-api/internal/features/users/repository/postgres"
)

func (r *AuthRepository) CreateUserWithIdentity(
	ctx context.Context,
	user domain.User,
	provider, subject, email string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	userQuery := `
		insert into todoapp.users (full_name, phone_number)
		values($1, $2)
		returning id, version, full_name, phone_number;
	`

	identityQuery := `
		insert into todoapp.user_identities(user_id, provider, subject, email)
		values ($1, $2, $3, $4)
		returning id, user_id, provider, subject, email, created_at;
	`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// user tx
	row := tx.QueryRow(ctx, userQuery, user.FullName, user.PhoneNumber)

	var userModel users_postgres_repository.UserModel
	err = row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("user scan error: %w", err)
	}

	// identity tx
	_, err = tx.Exec(ctx, identityQuery, userModel.ID, provider, subject, email)
	if err != nil {
		return domain.User{}, fmt.Errorf("identity scan error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("commit tx: %w", err)
	}

	return domain.User{
		ID:          userModel.ID,
		Version:     userModel.Version,
		FullName:    userModel.FullName,
		PhoneNumber: userModel.PhoneNumber,
	}, nil
}
