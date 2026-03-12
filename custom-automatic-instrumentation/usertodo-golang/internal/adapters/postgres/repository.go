package postgres

import (
	"context"
	"errors"

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

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).Preload("Tasks").First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	user := &domain.User{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Tasks:     mapTasks(model.Tasks),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	return user, nil
}

func (r *UserRepository) AttachTask(ctx context.Context, userID uint, taskID uint) error {
	link := UserTaskModel{
		UserID: userID,
		TaskID: taskID,
	}

	return r.db.WithContext(ctx).Where("user_id = ? AND task_id = ?", userID, taskID).FirstOrCreate(&link).Error
}

func (r *TaskRepository) FindByIDs(ctx context.Context, ids []uint) ([]domain.Task, error) {
	var models []TaskModel

	if err := r.db.WithContext(ctx).Where("id IN ? AND done = ?", ids, false).Find(&models).Error; err != nil {
		return nil, err
	}

	return mapTasks(models), nil
}

func (r *TaskRepository) List(ctx context.Context) ([]domain.Task, error) {
	var models []TaskModel

	if err := r.db.WithContext(ctx).Where("done = ?", false).Order("created_at asc").Find(&models).Error; err != nil {
		return nil, err
	}

	return mapTasks(models), nil
}

func mapTasks(models []TaskModel) []domain.Task {
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

	return tasks
}
