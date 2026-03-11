package postgres

import (
	"context"

	"interview/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type TaskRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&UserModel{}, &UserTaskModel{})
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := UserModel{
		Name:  user.Name,
		Email: user.Email,
		Tasks: make([]TaskModel, 0, len(user.Tasks)),
	}

	for _, task := range user.Tasks {
		model.Tasks = append(model.Tasks, TaskModel{ID: task.ID})
	}

	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		user.ID = model.ID
		user.CreatedAt = model.CreatedAt
		user.UpdatedAt = model.UpdatedAt

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	model := TaskModel{
		Title:       task.Title,
		Description: task.Description,
		Done:        task.Done,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	task.ID = model.ID
	task.Done = model.Done
	task.CreatedAt = model.CreatedAt
	task.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *TaskRepository) FindByIDs(ctx context.Context, ids []uint) ([]domain.Task, error) {
	var models []TaskModel

	if err := r.db.WithContext(ctx).Where("id IN ? AND done = ?", ids, false).Find(&models).Error; err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, 0, len(models))
	for _, model := range models {
		tasks = append(tasks, domain.Task{
			ID:          model.ID,
			Title:       model.Title,
			Description: model.Description,
			Done:        model.Done,
			CreatedAt:   model.CreatedAt,
			UpdatedAt:   model.UpdatedAt,
		})
	}

	return tasks, nil
}
