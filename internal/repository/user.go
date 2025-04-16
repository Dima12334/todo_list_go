package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"todo_list_go/internal/models"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}

func (r *UserRepo) Update(ctx context.Context, inp models.UpdateUserInput) (*models.User, error) {
	return nil, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
