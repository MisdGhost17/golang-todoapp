package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
)

func (s *StatisticsService) GetTasks(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"to must be after from: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{
			TasksCreated:               0,
			TasksCompleted:             0,
			TasksCompletedRate:         nil,
			TasksAverageComplitionTime: nil,
		}
	}

	tasksCreated := len(tasks)

	tasksComleted := 0

	var totalCompletedTasksDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksComleted += 1
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletedTasksDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksComleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksComleted > 0 && totalCompletedTasksDuration != 0 {
		avg := totalCompletedTasksDuration / time.Duration(tasksComleted)

		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksComleted,
		&tasksCompletedRate,
		tasksAverageCompletionTime,
	)

}
