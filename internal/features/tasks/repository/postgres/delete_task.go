package tasks_postgres_repository

import (
	"context"
	"fmt"
	core_errors "todo-api/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		delete
		from todoapp.tasks
		where id = $1
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf(
			"delete task with id=%d: %w",
			id,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
