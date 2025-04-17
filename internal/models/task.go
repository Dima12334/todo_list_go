package models

import "time"

type Task struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
}

type TaskCategory struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	CategoryID string `json:"category_id"`
}

type UpdateTaskInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Completed   bool     `json:"completed"`
	CategoryIDs []string `json:"category_ids"`
}
