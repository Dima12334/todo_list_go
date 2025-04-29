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

func (s *CategoryService) Create(ctx context.Context, inp CreateCategoryInput) (models.Category, error) {
	category := models.Category{
		UserID:      inp.UserID,
		Title:       inp.Title,
		Description: inp.Description,
		Color:       inp.Color,
	}
	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		return models.Category{}, err
	}
	return createdCategory, nil
}

func (s *CategoryService) Update(ctx context.Context, inp UpdateCategoryInput) (models.Category, error) {
	_, err := s.repo.GetByID(ctx, inp.ID, inp.UserID)
	if err != nil {
		return models.Category{}, err
	}

	updateInput := repository.UpdateCategoryInput{
		ID:          inp.ID,
		Title:       inp.Title,
		Description: inp.Description,
		Color:       inp.Color,
	}

	return s.repo.Update(ctx, updateInput)
}

func (s *CategoryService) Delete(ctx context.Context, CategoryID, UserID string) error {
	_, err := s.repo.GetByID(ctx, CategoryID, UserID)
	if err != nil {
		return err
	}

	return s.repo.DeleteByID(ctx, CategoryID)
}

func (s *CategoryService) GetList(ctx context.Context, userID string) ([]models.Category, error) {
	return s.repo.GetListByUserID(ctx, userID)
}
