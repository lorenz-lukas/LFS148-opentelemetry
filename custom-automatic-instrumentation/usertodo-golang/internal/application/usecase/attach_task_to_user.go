package usecase

import (
	"context"
	"errors"

	"interview/internal/application/ports"
	"interview/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type AttachTaskToUserInput struct {
	UserID uint
	TaskID uint
}

type AttachTaskToUserUseCase struct {
	userRepository ports.UserRepository
	taskRepository ports.TaskRepository
}

func NewAttachTaskToUserUseCase(userRepository ports.UserRepository, taskRepository ports.TaskRepository) *AttachTaskToUserUseCase {
	return &AttachTaskToUserUseCase{
		userRepository: userRepository,
		taskRepository: taskRepository,
	}
}

func (uc *AttachTaskToUserUseCase) Execute(ctx context.Context, input AttachTaskToUserInput) (*domain.User, error) {
	user, err := uc.userRepository.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	tasks, err := uc.taskRepository.FindByIDs(ctx, []uint{input.TaskID})
	if err != nil {
		return nil, err
	}
	if len(tasks) != 1 {
		return nil, ErrTaskNotFound
	}

	if err := uc.userRepository.AttachTask(ctx, input.UserID, input.TaskID); err != nil {
		return nil, err
	}

	return uc.userRepository.FindByID(ctx, input.UserID)
}
