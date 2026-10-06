package users_postgres_repository

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		insert into todoapp.users (full_name, phone_number)
		values ($1, $2)
		returning id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.User{
		ID:          userModel.ID,
		Version:     userModel.Version,
		FullName:    userModel.FullName,
		PhoneNumber: userModel.PhoneNumber,
	}, nil
}
