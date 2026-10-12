package tasks_transport_http

import (
	"time"
	"todo-api/internal/core/domain"
)

type TaskDTOResponse struct {
	ID      int `json:"id"`
	Version int `json:"version"`

	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`

	AuthorUserID int `json:"authorId"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func tasksDTOFromDomains(tasks []domain.Task) []TaskDTOResponse {
	tasksDTOs := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		tasksDTOs[i] = taskDTOFromDomain(task)
	}

	return tasksDTOs
}
