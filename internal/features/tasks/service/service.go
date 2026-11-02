package tasks_service

import (
	"context"
	"todo-api/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		limit, offset *int,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		id int,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id int,
	) error

	TitleExists(
		ctx context.Context,
		task domain.Task,
	) (bool, error)
}

func NewTasksService(repo TasksRepository) *TasksService {
	return &TasksService{tasksRepository: repo}
}
