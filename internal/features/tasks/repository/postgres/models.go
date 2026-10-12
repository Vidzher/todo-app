package tasks_postgres_repository

import (
	"time"
	"todo-api/internal/core/domain"
)

type TaskModel struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

func tasksDomainsFormModels(tasks []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(tasks))

	for i, task := range tasks {
		domains[i] = domain.NewTask(
			task.ID,
			task.Version,
			task.AuthorUserID,
			task.Title,
			task.Description,
			task.Completed,
			task.CreatedAt,
			task.CompletedAt,
		)
	}

	return domains
}
