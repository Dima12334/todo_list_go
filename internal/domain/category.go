package domain

import "time"

type Category struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Color       string    `json:"color" db:"color"`
}
