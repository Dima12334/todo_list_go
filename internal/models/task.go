package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
}

type TaskCategory struct {
	ID         int `json:"id"`
	TaskID     int `json:"task_id"`
	CategoryID int `json:"category_id"`
}

type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CategoryIDs []int  `json:"category_ids"`
}
