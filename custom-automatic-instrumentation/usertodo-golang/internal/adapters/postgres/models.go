package postgres

import "time"

type UserModel struct {
	ID        uint        `gorm:"primaryKey"`
	Name      string      `gorm:"size:120;not null"`
	Email     string      `gorm:"size:255;not null;uniqueIndex"`
	Tasks     []TaskModel `gorm:"many2many:user_tasks;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TaskModel struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:120;not null"`
	Description string `gorm:"column:description;size:255;not null;default:''"`
	Done        bool   `gorm:"column:done;not null;default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserTaskModel struct {
	UserID uint `gorm:"primaryKey"`
	TaskID uint `gorm:"primaryKey"`
}

func (UserModel) TableName() string {
	return "users"
}

func (TaskModel) TableName() string {
	return "todo"
}

func (UserTaskModel) TableName() string {
	return "user_tasks"
}
