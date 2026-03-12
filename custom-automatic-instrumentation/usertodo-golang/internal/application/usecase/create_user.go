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
	Name  string
	Email string
}

type CreateUserUseCase struct {
	userRepository ports.UserRepository
}

func NewCreateUserUseCase(userRepository ports.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
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

	if err := uc.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
