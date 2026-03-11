package ports

import (
	"context"

	"interview/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
}

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	FindByIDs(ctx context.Context, ids []uint) ([]domain.Task, error)
}
