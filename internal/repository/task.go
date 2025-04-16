package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"todo_list_go/internal/models"
)

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	return nil, nil
}

func (r *TaskRepo) Update(ctx context.Context, inp models.UpdateTaskInput) (*models.Task, error) {
	return nil, nil
}

func (r *TaskRepo) DeleteByID(ctx context.Context, id string) error {
	return nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	return nil, nil
}

func (r *TaskRepo) GetListByUserID(ctx context.Context, userID string) ([]*models.Task, error) {
	return nil, nil
}

func (r *TaskRepo) GetListByCategoryIDs(ctx context.Context, categoryIDs []string) ([]*models.Task, error) {
	return nil, nil
}
