package tasks_postgres_repository

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
)

func (r *TasksRepository) TitleExists(
	ctx context.Context,
	task domain.Task,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `select exists(
		select 1
		from todoapp.tasks
		where author_user_id = $1 and title = $2
	)`

	row := r.pool.QueryRow(ctx, query, task.AuthorUserID, task.Title)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("scan error: %w", err)
	}

	return exists, nil
}
