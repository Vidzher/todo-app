package tasks_service

import (
	"context"
	"fmt"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	// move validation to transport layer
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	exists, err := s.tasksRepository.TitleExists(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	if exists {
		return domain.Task{}, fmt.Errorf(
			"task with title=%s already exists: %w",
			task.Title,
			core_errors.ErrConflict,
		)
	}

	task, err = s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
