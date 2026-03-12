package usecase

import (
	"context"

	"interview/internal/application/ports"
	"interview/internal/domain"
)

type ListTasksUseCase struct {
	taskRepository ports.TaskRepository
}

func NewListTasksUseCase(taskRepository ports.TaskRepository) *ListTasksUseCase {
	return &ListTasksUseCase{taskRepository: taskRepository}
}

func (uc *ListTasksUseCase) Execute(ctx context.Context) ([]domain.Task, error) {
	return uc.taskRepository.List(ctx)
}
