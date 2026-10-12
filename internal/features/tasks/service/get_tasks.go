package tasks_service

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	limit, offset *int,
) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be greater than or equal to 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be greater or equal to 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}

	return tasks, nil
}
