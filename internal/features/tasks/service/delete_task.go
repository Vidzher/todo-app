package tasks_service

import (
	"context"
	"fmt"
	core_errors "todo-api/internal/core/errors"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	id int,
) error {
	if id <= 0 {
		return fmt.Errorf(
			"taskId must be greater oq equal to 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if err := s.tasksRepository.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
