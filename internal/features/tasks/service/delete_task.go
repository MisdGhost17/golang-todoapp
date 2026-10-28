package tasks_service

import (
	"context"
	"fmt"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	taskID int,
	userID int,
) error {
	if err := s.taskRepository.DeleteTask(ctx, taskID, userID); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}
	return nil
}
