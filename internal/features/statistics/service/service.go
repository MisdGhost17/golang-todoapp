package statistics_service

import (
	"context"
	"time"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}
