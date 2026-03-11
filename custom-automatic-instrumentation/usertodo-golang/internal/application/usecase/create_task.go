package usecase

import (
	"context"
	"fmt"
	"strings"

	"interview/internal/application/ports"
	"interview/internal/domain"
)

type CreateTaskInput struct {
	Title       string
	Description string
}

type CreateTaskUseCase struct {
	taskRepository ports.TaskRepository
}

func NewCreateTaskUseCase(taskRepository ports.TaskRepository) *CreateTaskUseCase {
	return &CreateTaskUseCase{taskRepository: taskRepository}
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, input CreateTaskInput) (*domain.Task, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}

	task := &domain.Task{
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Done:        false,
	}

	if err := uc.taskRepository.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}
