package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"interview/internal/application/ports"
	"interview/internal/domain"
)

var ErrTaskNotFound = errors.New("one or more tasks were not found")

type CreateUserInput struct {
	Name    string
	Email   string
	TaskIDs []uint
}

type CreateUserUseCase struct {
	userRepository ports.UserRepository
	taskRepository ports.TaskRepository
}

func NewCreateUserUseCase(userRepository ports.UserRepository, taskRepository ports.TaskRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
		taskRepository: taskRepository,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, fmt.Errorf("name is required")
	}

	if strings.TrimSpace(input.Email) == "" {
		return nil, fmt.Errorf("email is required")
	}

	user := &domain.User{
		Name:  strings.TrimSpace(input.Name),
		Email: strings.TrimSpace(input.Email),
	}

	if len(input.TaskIDs) > 0 {
		tasks, err := uc.taskRepository.FindByIDs(ctx, input.TaskIDs)
		if err != nil {
			return nil, err
		}

		if len(tasks) != len(uniqueUintSlice(input.TaskIDs)) {
			return nil, ErrTaskNotFound
		}

		user.Tasks = tasks
	}

	if err := uc.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func uniqueUintSlice(values []uint) []uint {
	seen := make(map[uint]struct{}, len(values))
	result := make([]uint, 0, len(values))

	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}
