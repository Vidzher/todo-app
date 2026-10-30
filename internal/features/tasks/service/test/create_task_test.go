package tasks_service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"
	tasks_service "todo-api/internal/features/tasks/service"
)

func TestCreateTask(t *testing.T) {
	cases := []struct {
		name        string
		input       domain.Task
		expected    domain.Task
		expectedErr error
		exists      bool
	}{
		{
			name: "VALID: Create task",
			input: domain.Task{
				ID:           domain.UninitializedID,
				Version:      domain.UninitializedVersion,
				AuthorUserID: 1,
				Title:        "New task",
				Description:  nil,
				Completed:    false,
				CreatedAt:    time.Now(),
				CompletedAt:  nil,
			},
			expected: domain.Task{
				ID:           1,
				Version:      1,
				AuthorUserID: 1,
				Title:        "New task",
				Description:  nil,
				Completed:    false,
				CreatedAt:    time.Now(),
				CompletedAt:  nil,
			},
			expectedErr: nil,
			exists:      false,
		},
		{
			name: "INVALID: Create task",
			input: domain.Task{
				ID:           domain.UninitializedID,
				Version:      domain.UninitializedVersion,
				AuthorUserID: 1,
				Description:  nil,
				Completed:    false,
				CreatedAt:    time.Now(),
				CompletedAt:  nil,
			},
			expected:    domain.Task{},
			expectedErr: core_errors.ErrInvalidArgument,
			exists:      false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := &mockTaskRepository{
				task:   c.expected,
				err:    c.expectedErr,
				exists: c.exists,
			}

			service := tasks_service.NewTasksService(repo)

			task, err := service.CreateTask(context.Background(), c.input)

			if !errors.Is(err, c.expectedErr) {
				t.Fatalf("error = %v, want %v", err, c.expectedErr)
			}

			if c.expectedErr == nil &&
				!reflect.DeepEqual(task, c.expected) {
				t.Fatalf("result = %+v, want %+v", task, c.expected)
			}
		})
	}
}

type mockTaskRepository struct {
	task   domain.Task
	err    error
	exists bool
}

func (r *mockTaskRepository) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	return r.task, r.err
}

func (r *mockTaskRepository) TitleExists(
	ctx context.Context,
	task domain.Task,
) (bool, error) {
	return r.exists, r.err
}

func (r *mockTaskRepository) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	panic("not implemented")
}

func (r *mockTaskRepository) GetTasks(
	ctx context.Context,
	limit, offset *int,
) ([]domain.Task, error) {
	panic("not implemented")
}

func (r *mockTaskRepository) DeleteTask(
	ctx context.Context,
	id int,
) error {
	panic("not implemented")
}
