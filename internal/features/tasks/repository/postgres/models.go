package tasks_postgres_repository

import (
	"time"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
)

type TaskModel struct {
	ID           int        `postgres:"id"`
	Version      int        `postgres:"version"`
	Title        string     `postgres:"title"`
	Description  *string    `postgres:"description"`
	Completed    bool       `postgres:"completed"`
	CreatedAt    time.Time  `postgres:"created_at"`
	CompletedAt  *time.Time `postgres:"completed_at"`
	AuthorUserID int        `postgres:"author_user_id"`
}

func taskDomainsFromModels(taskModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(taskModels))

	for i, task := range taskModels {
		domains[i] = taskDomainFromModel(task)
	}

	return domains
}

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}
