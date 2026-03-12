package ports

import (
	"context"

	"interview/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	AttachTask(ctx context.Context, userID uint, taskID uint) error
}

type TaskRepository interface {
	FindByIDs(ctx context.Context, ids []uint) ([]domain.Task, error)
	List(ctx context.Context) ([]domain.Task, error)
}
