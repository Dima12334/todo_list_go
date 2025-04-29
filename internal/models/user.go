package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Password  string    `json:"password" db:"password"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
}

type UpdateUserInput struct {
	Name string `json:"name"`
}
