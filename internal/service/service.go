package service

import (
	"context"
	"time"
	"todo_list_go/internal/domain"
	"todo_list_go/internal/repository"
	"todo_list_go/pkg/auth"
	"todo_list_go/pkg/hash"
)

//go:generate mockgen -source=service.go -destination=mocks/mock_service.go

type SignUpUserInput struct {
	Name     string
	Email    string
	Password string
}

type SignInUserInput struct {
	Email    string
	Password string
}

type User interface {
	SignUp(ctx context.Context, inp SignUpUserInput) error
	SignIn(ctx context.Context, inp SignInUserInput) (string, error)
	GetByID(ctx context.Context, userID string) (domain.User, error)
}

type CreateTaskInput struct {
	UserID      string `json:"user_id"`
	CategoryID  string `json:"category_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type UpdateTaskInput struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	CategoryID  *string `json:"category_id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

type TaskOutput struct {
	ID          string          `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Category    domain.Category `json:"category"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Completed   bool            `json:"completed"`
}

type TaskListResult struct {
	Items      []TaskOutput
	TotalItems int64
	TotalPages int
}

type Task interface {
	Create(ctx context.Context, inp CreateTaskInput) (TaskOutput, error)
	Update(ctx context.Context, inp UpdateTaskInput) (TaskOutput, error)
	Delete(ctx context.Context, taskID, userID string) error
	GetByID(ctx context.Context, taskID, userID string) (TaskOutput, error)
	GetList(ctx context.Context, userID string, query domain.GetTasksQuery) (TaskListResult, error)
}

type CreateCategoryInput struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateCategoryInput struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

type Category interface {
	Create(ctx context.Context, inp CreateCategoryInput) (domain.Category, error)
	Update(ctx context.Context, inp UpdateCategoryInput) (domain.Category, error)
	Delete(ctx context.Context, categoryID, userID string) error
	GetList(ctx context.Context, userID string) ([]domain.Category, error)
}

type Deps struct {
	Repos          *repository.Repositories
	AccessTokenTTL time.Duration
	TokenManager   auth.TokenManager
	Hasher         hash.PasswordHasher
}

type Services struct {
	Users      User
	Tasks      Task
	Categories Category
}

func NewServices(deps Deps) *Services {
	return &Services{
		Users:      NewUserService(deps.Repos.User, deps.AccessTokenTTL, deps.TokenManager, deps.Hasher),
		Tasks:      NewTaskService(deps.Repos.Task, deps.Repos.Category),
		Categories: NewCategoryService(deps.Repos.Category),
	}
}
