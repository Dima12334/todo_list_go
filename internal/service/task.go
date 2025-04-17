package service

import (
	"context"
	"todo_list_go/internal/models"
	"todo_list_go/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(ctx context.Context, inp CreateTaskInput) (*models.Task, error) {
	return nil, nil
}

func (s *TaskService) Update(ctx context.Context, inp UpdateTaskInput) (*models.Task, error) {
	return nil, nil
}

func (s *TaskService) Delete(ctx context.Context, TaskID, UserID string) error {
	return nil
}

func (s *TaskService) GetByID(ctx context.Context, TaskID, UserID string) (*models.Task, error) {
	return nil, nil
}

func (s *TaskService) GetList(ctx context.Context, userID string, categoryIDs []string) ([]*models.Task, error) {
	return nil, nil
}
