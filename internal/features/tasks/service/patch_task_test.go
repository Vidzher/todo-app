package tasks_service

import (
	"context"
	"testing"
	"time"
	"todo-api/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	createdAtMock = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	taskID        = 1
	authorUserID  = 1
)

func TestPatchTask(t *testing.T) {
	cases := []struct {
		name              string
		initial           domain.Task
		input             domain.TaskPatch
		expectedCompleted bool
	}{
		{
			name: "VALID: Mark as completed",
			initial: domain.Task{
				ID:           taskID,
				Version:      1,
				Title:        "Prev",
				Description:  new("Prev"),
				Completed:    false,
				CreatedAt:    createdAtMock,
				CompletedAt:  nil,
				AuthorUserID: authorUserID,
			},
			input: domain.TaskPatch{
				Completed:   new(true),
				Title:       nil,
				Description: nil,
			},
			expectedCompleted: true,
		},
		{
			name: "VALID: Mark as uncompleted",
			initial: domain.Task{
				ID:           taskID,
				Version:      1,
				Title:        "Prev",
				Description:  new("Prev"),
				Completed:    true,
				CreatedAt:    createdAtMock,
				CompletedAt:  new(createdAtMock.Add(time.Hour)),
				AuthorUserID: authorUserID,
			},
			input: domain.TaskPatch{
				Completed:   new(false),
				Title:       nil,
				Description: nil,
			},
			expectedCompleted: false,
		},
		{
			name: "VALID: Change title and description",
			initial: domain.Task{
				ID:           taskID,
				Version:      1,
				Title:        "Prev",
				Description:  new("Prev"),
				Completed:    true,
				CreatedAt:    createdAtMock,
				CompletedAt:  nil,
				AuthorUserID: authorUserID,
			},
			input: domain.TaskPatch{
				Completed:   new(false),
				Title:       new("Next"),
				Description: new("Next"),
			},
			expectedCompleted: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := NewMockTasksRepository(t)

			repo.EXPECT().
				GetTask(context.Background(), taskID).
				Return(c.initial, nil).
				Once()

			repo.EXPECT().
				PatchTask(context.Background(), taskID, mock.Anything).
				RunAndReturn(func(_ context.Context, _ int, task domain.Task) (domain.Task, error) {
					return domain.Task{
						ID:           task.ID,
						Version:      task.Version + 1,
						Title:        task.Title,
						Description:  task.Description,
						Completed:    task.Completed,
						CreatedAt:    task.CreatedAt,
						CompletedAt:  task.CompletedAt,
						AuthorUserID: task.AuthorUserID,
					}, nil
				}).
				Once()

			service := NewTasksService(repo)

			task, err := service.PatchTask(context.Background(), taskID, c.input)

			require.NoError(t, err)
			assert.Equal(t, c.initial.ID, task.ID)
			assert.False(t, c.initial.Version == task.Version)

			assert.Equal(t, c.expectedCompleted, task.Completed)
			assert.Equal(t, c.initial.CreatedAt, task.CreatedAt)

			if c.input.Title != nil {
				assert.True(t, c.initial.Title != task.Title)
			} else {
				assert.Equal(t, c.initial.Title, task.Title)
			}

			if c.input.Description != nil {
				assert.True(t, c.initial.Description != task.Description)
			} else {
				assert.Equal(t, c.initial.Description, task.Description)
			}

			if c.expectedCompleted {
				require.NotNil(t, task.CompletedAt)
				assert.False(t, task.CompletedAt.Before(c.initial.CreatedAt))
			} else {
				assert.Nil(t, task.CompletedAt)
			}
		})
	}
}
