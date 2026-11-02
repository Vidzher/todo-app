package tasks_service

import (
	"context"
	"testing"
	"time"
	"todo-api/internal/core/domain"
	core_errors "todo-api/internal/core/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	cases := []struct {
		name              string
		input             domain.Task
		expected          domain.Task
		expectedErr       error
		exists            bool
		wrongAuthorUserId bool
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
			expectedErr:       nil,
			exists:            false,
			wrongAuthorUserId: false,
		},
		{
			name: "INVALID: Without title",
			input: domain.Task{
				ID:           domain.UninitializedID,
				Version:      domain.UninitializedVersion,
				AuthorUserID: 1,
				Description:  nil,
				Completed:    false,
				CreatedAt:    time.Now(),
				CompletedAt:  nil,
			},
			expected:          domain.Task{},
			expectedErr:       core_errors.ErrInvalidArgument,
			exists:            false,
			wrongAuthorUserId: false,
		},
		{
			name: "INVALID: With collision",
			input: domain.Task{
				ID:           domain.UninitializedID,
				Version:      domain.UninitializedVersion,
				AuthorUserID: 1,
				Title:        "Task title",
				Description:  nil,
				Completed:    false,
				CreatedAt:    time.Now(),
				CompletedAt:  nil,
			},
			expected:          domain.Task{},
			expectedErr:       core_errors.ErrConflict,
			exists:            true,
			wrongAuthorUserId: false,
		},
		{
			name: "INVALID: Wrong AuthorUserID",
			input: domain.Task{
				ID:           domain.UninitializedID,
				Version:      domain.UninitializedVersion,
				AuthorUserID: -1,
				Title:        "Task title",
				Description:  nil,
				Completed:    false,
				CreatedAt:    time.Now(),
				CompletedAt:  nil,
			},
			expected:          domain.Task{},
			expectedErr:       core_errors.ErrInvalidArgument,
			exists:            false,
			wrongAuthorUserId: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := NewMockTasksRepository(t)

			if c.input.Title != "" && !c.wrongAuthorUserId {
				repo.EXPECT().
					TitleExists(context.Background(), c.input).
					Return(c.exists, c.expectedErr).
					Once()
			}

			if c.expectedErr == nil && !c.exists {
				repo.EXPECT().
					CreateTask(context.Background(), c.input).
					Return(c.expected, nil).
					Once()
			}

			service := NewTasksService(repo)

			task, err := service.CreateTask(context.Background(), c.input)

			if c.expectedErr != nil {
				require.ErrorIs(t, err, c.expectedErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.expected, task)
		})
	}
}
