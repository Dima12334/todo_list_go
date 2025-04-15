package models

import "time"

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}
