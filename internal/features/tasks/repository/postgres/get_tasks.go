package tasks_postgres_repository

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	limit, offset *int,
) ([]domain.Task, error) {
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
		order by id
		limit $1
		offset $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel
		err := rows.Scan(
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
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	tasks := tasksDomainsFormModels(taskModels)

	return tasks, nil
}
