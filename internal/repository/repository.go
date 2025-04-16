package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"todo_list_go/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, inp models.UpdateUserInput) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	Update(ctx context.Context, inp models.UpdateTaskInput) (*models.Task, error)
	DeleteByID(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	GetListByUserID(ctx context.Context, userID string) ([]*models.Task, error)
	GetListByCategoryIDs(ctx context.Context, categoryIDs []string) ([]*models.Task, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) (*models.Category, error)
	Update(ctx context.Context, inp models.UpdateCategoryInput) (*models.Category, error)
	DeleteByID(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.Category, error)
	GetListByUserID(ctx context.Context, userID string) ([]*models.Category, error)
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
