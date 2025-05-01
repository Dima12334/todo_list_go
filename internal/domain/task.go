package domain

import "time"

type Task struct {
	ID          string    `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UserID      string    `json:"user_id" db:"user_id"`
	CategoryID  string    `json:"category_id" db:"category_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Completed   bool      `json:"completed" db:"completed"`
}
