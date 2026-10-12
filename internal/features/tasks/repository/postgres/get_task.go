package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
	core_postgres_pool "todo-api/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		select
			id,
			version,
			title,
			description,
			completed,
			created_at,
			completed_at,
			author_user_id
		from todoapp.tasks
		where id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

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
	if errors.Is(err, core_postgres_pool.ErrNoRows) {
		return domain.Task{}, fmt.Errorf(
			"task with id=%d: %v: %w",
			id,
			err,
			core_errors.ErrNotFound,
		)
	}

	if err != nil {
		return domain.Task{}, fmt.Errorf("scan task: %w", err)
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
