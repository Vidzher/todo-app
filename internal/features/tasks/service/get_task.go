package tasks_service

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	if id <= 0 {
		return domain.Task{}, fmt.Errorf(
			"id must be greater than 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	return task, nil
}
