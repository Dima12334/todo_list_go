package service

import (
	"context"
	"time"
	"todo_list_go/internal/models"
	"todo_list_go/internal/repository"
	customErrors "todo_list_go/pkg/errors"
)

type TaskService struct {
	repo         repository.TaskRepository
	categoryRepo repository.CategoryRepository
}

func NewTaskService(repo repository.TaskRepository, categoryRepo repository.CategoryRepository) *TaskService {
	return &TaskService{repo: repo, categoryRepo: categoryRepo}
}

func (s *TaskService) Create(ctx context.Context, inp CreateTaskInput) (TaskOutput, error) {
	_, err := s.categoryRepo.GetByID(ctx, inp.CategoryID, inp.UserID)
	if err != nil {
		return TaskOutput{}, err
	}

	task := models.Task{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      inp.UserID,
		CategoryID:  inp.CategoryID,
		Title:       inp.Title,
		Description: inp.Description,
		Completed:   inp.Completed,
	}
	createdTask, err := s.repo.Create(ctx, task)
	if err != nil {
		return TaskOutput{}, err
	}

	return TaskOutput(createdTask), nil
}

func (s *TaskService) Update(ctx context.Context, inp UpdateTaskInput) (TaskOutput, error) {
	_, err := s.repo.GetByID(ctx, inp.ID, inp.UserID)
	if err != nil {
		return TaskOutput{}, err
	}

	if inp.Title == nil && inp.Description == nil && inp.CategoryID == nil && inp.Completed == nil {
		return TaskOutput{}, customErrors.ErrNoUpdateFields
	}

	if inp.CategoryID != nil {
		_, err = s.categoryRepo.GetByID(ctx, *inp.CategoryID, inp.UserID)
		if err != nil {
			return TaskOutput{}, err
		}
	}

	updateInput := repository.UpdateTaskInput{
		ID:          inp.ID,
		UserID:      inp.UserID,
		UpdatedAt:   time.Now(),
		CategoryID:  inp.CategoryID,
		Title:       inp.Title,
		Description: inp.Description,
		Completed:   inp.Completed,
	}
	updatedTask, err := s.repo.Update(ctx, updateInput)
	if err != nil {
		return TaskOutput{}, err
	}

	return TaskOutput(updatedTask), nil
}

func (s *TaskService) Delete(ctx context.Context, taskID, userID string) error {
	_, err := s.repo.GetByID(ctx, taskID, userID)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, taskID)
}

func (s *TaskService) GetByID(ctx context.Context, taskID, userID string) (TaskOutput, error) {
	task, err := s.repo.GetByID(ctx, taskID, userID)
	if err != nil {
		return TaskOutput{}, err
	}

	return TaskOutput(task), nil
}

func (s *TaskService) GetList(ctx context.Context, userID string) ([]TaskOutput, error) {
	tasks, err := s.repo.GetListByUserID(ctx, userID)
	if err != nil {
		return []TaskOutput{}, err
	}

	tasksOutput := make([]TaskOutput, len(tasks))
	for i, task := range tasks {
		tasksOutput[i] = TaskOutput(task)
	}
	return tasksOutput, nil
}
