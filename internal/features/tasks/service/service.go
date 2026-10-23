package tasks_service

import (
	"context"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
)

type TasksService struct {
	taskRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		userID *int,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		taskID int,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		taskID int,
	) error

	PatchTask(
		ctx context.Context,
		taskID int,
		task domain.Task,
	) (domain.Task, error)
}

func NewTasksService(
	taskRepository TasksRepository,
) *TasksService {
	return &TasksService{
		taskRepository: taskRepository,
	}
}
