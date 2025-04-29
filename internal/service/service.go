package service

import (
	"context"
	"time"
	"todo_list_go/internal/models"
	"todo_list_go/internal/repository"
	"todo_list_go/pkg/auth"
	"todo_list_go/pkg/hash"
)

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
	GetByID(ctx context.Context, userID string) (models.User, error)
}

type CreateTaskInput struct {
	UserID      string   `json:"user_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	CategoryIDs []string `json:"category_ids"`
}

type UpdateTaskInput struct {
	TaskID      string   `json:"task_id"`
	UserID      string   `json:"user_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Completed   bool     `json:"completed"`
	CategoryIDs []string `json:"category_ids"`
}

type Task interface {
	Create(ctx context.Context, inp CreateTaskInput) (models.Task, error)
	Update(ctx context.Context, inp UpdateTaskInput) (models.Task, error)
	Delete(ctx context.Context, TaskID, UserID string) error
	GetByID(ctx context.Context, TaskID, UserID string) (models.Task, error)
	GetList(ctx context.Context, userID string, categoryIDs []string) ([]models.Task, error)
}

type CreateCategoryInput struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateCategoryInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type Category interface {
	Create(ctx context.Context, inp CreateCategoryInput) (models.Category, error)
	Update(ctx context.Context, inp UpdateCategoryInput) (models.Category, error)
	Delete(ctx context.Context, CategoryID, UserID string) error
	GetList(ctx context.Context, userID string) ([]models.Category, error)
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
		Tasks:      NewTaskService(deps.Repos.Task),
		Categories: NewCategoryService(deps.Repos.Category),
	}
}
