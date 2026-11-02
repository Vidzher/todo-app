package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
	core_postgres_pool "todo-api/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	id int,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		update todoapp.tasks
		set
			title=$1,
			description=$2,
			completed=$3,
			completed_at=$4,
			version=version+1
		where id=$5 and version=$6
		returning
			id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	row := r.pool.QueryRow(
		ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		id,
		task.Version,
	)

	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id=%d concurrently accessed: %w",
				id,
				core_errors.ErrConflict,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.Task{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CompletedAt:  task.CompletedAt,
		CreatedAt:    task.CreatedAt,
		AuthorUserID: task.AuthorUserID,
	}, nil
}
