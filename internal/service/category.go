package service

import (
	"context"
	"todo_list_go/internal/models"
	"todo_list_go/internal/repository"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, inp CreateCategoryInput) (*models.Category, error) {
	return nil, nil
}

func (s *CategoryService) Update(ctx context.Context, inp UpdateCategoryInput) (*models.Category, error) {
	return nil, nil
}

func (s *CategoryService) Delete(ctx context.Context, CategoryID, UserID string) error {
	return nil
}

func (s *CategoryService) GetList(ctx context.Context, userID string) ([]*models.Category, error) {
	return nil, nil
}
