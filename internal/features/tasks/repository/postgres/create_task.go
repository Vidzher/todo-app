package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
	core_postgres_pool "todo-api/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		insert into todoapp.tasks (title, description, author_user_id)
		values ($1, $2, $3)
		returning 
			id,
			version,
			title, 
			description, 
			completed, 
			created_at, 
			completed_at,
			author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.AuthorUserID,
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
		if errors.Is(err, core_postgres_pool.ErrViolatesFK) {
			return domain.Task{}, fmt.Errorf(
				"authorID=%d not found: %v: %w",
				task.AuthorUserID,
				err,
				core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.Task{
		ID:           taskModel.ID,
		Version:      taskModel.Version,
		Title:        taskModel.Title,
		Description:  taskModel.Description,
		Completed:    taskModel.Completed,
		CreatedAt:    taskModel.CreatedAt,
		CompletedAt:  taskModel.CompletedAt,
		AuthorUserID: taskModel.AuthorUserID,
	}, nil
}
