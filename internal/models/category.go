package models

import "time"

type Category struct {
	ID          int       `json:"id"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
}

type UpdateCategoryInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}
