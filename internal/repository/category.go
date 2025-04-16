package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"todo_list_go/internal/models"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, category *models.Category) (*models.Category, error) {
	return nil, nil
}

func (r *CategoryRepo) Update(ctx context.Context, inp models.UpdateCategoryInput) (*models.Category, error) {
	return nil, nil
}

func (r *CategoryRepo) DeleteByID(ctx context.Context, id string) error {
	return nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id string) (*models.Category, error) {
	return nil, nil
}

func (r *CategoryRepo) GetListByUserID(ctx context.Context, userID string) ([]*models.Category, error) {
	return nil, nil
}
