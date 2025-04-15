package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserId      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
}

type TaskCategory struct {
	ID         int `json:"id"`
	TaskId     int `json:"task_id"`
	CategoryId int `json:"category_id"`
}
