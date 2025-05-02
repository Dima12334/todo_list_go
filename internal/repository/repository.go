package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
	"todo_list_go/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type UpdateTaskInput struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	UpdatedAt   time.Time `json:"updated_at"`
	CategoryID  *string   `json:"category_id"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Completed   *bool     `json:"completed"`
}

type TaskOutput struct {
	ID          string          `json:"id" db:"id"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	Category    domain.Category `json:"category"`
	Title       string          `json:"title" db:"title"`
	Description string          `json:"description" db:"description"`
	Completed   bool            `json:"completed" db:"completed"`
}

type TaskRepository interface {
	Create(ctx context.Context, task domain.Task) (TaskOutput, error)
	Update(ctx context.Context, inp UpdateTaskInput) (TaskOutput, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, taskID, userID string) (TaskOutput, error)
	GetListByUserID(ctx context.Context, userID string, query domain.GetTasksQuery) ([]TaskOutput, int64, error)
}

type UpdateCategoryInput struct {
	ID          string  `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category domain.Category) (domain.Category, error)
	Update(ctx context.Context, inp UpdateCategoryInput) (domain.Category, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, categoryID, userID string) (domain.Category, error)
	GetListByUserID(ctx context.Context, userID string) ([]domain.Category, error)
}

type Repositories struct {
	User     UserRepository
	Task     TaskRepository
	Category CategoryRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		User:     NewUserRepo(db),
		Task:     NewTaskRepo(db),
		Category: NewCategoryRepo(db),
	}
}
