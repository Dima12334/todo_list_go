package service

import (
	"context"
	"time"
	"todo_list_go/internal/models"
	"todo_list_go/internal/repository"
	"todo_list_go/pkg/hash"
)

type UserService struct {
	repo           repository.UserRepository
	accessTokenTTL time.Duration
}

func NewUserService(repo repository.UserRepository, accessTokenTTL time.Duration) *UserService {
	return &UserService{
		repo:           repo,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *UserService) SignUp(ctx context.Context, inp SignUpUserInput) error {
	passwordHash, err := hash.GeneratePasswordHash(inp.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:      inp.Name,
		Password:  passwordHash,
		Email:     inp.Email,
		CreatedAt: time.Now(),
	}

	if err = s.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) SignIn(ctx context.Context, inp SignInUserInput) (string, error) {
	return "", nil
}

func (s *UserService) Update(ctx context.Context, inp UpdateUserInput) (*models.User, error) {
	return nil, nil
}
