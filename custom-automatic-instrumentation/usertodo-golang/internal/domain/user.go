package domain

import "time"

type User struct {
	ID        uint
	Name      string
	Email     string
	Tasks     []Task
	CreatedAt time.Time
	UpdatedAt time.Time
}
