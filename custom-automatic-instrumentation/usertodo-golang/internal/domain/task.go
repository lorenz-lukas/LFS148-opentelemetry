package domain

import "time"

type Task struct {
	ID          uint
	Title       string
	Description string
	Done        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
