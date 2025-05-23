package service

import (
	"context"
	"time"
	"todo_list_go/internal/domain"
	"todo_list_go/internal/repository"
	"todo_list_go/pkg/auth"
	customErrors "todo_list_go/pkg/errors"
	"todo_list_go/pkg/hash"
)

type UserService struct {
	repo           repository.UserRepository
	accessTokenTTL time.Duration
	tokenManager   auth.TokenManager
	hasher         hash.PasswordHasher
}

func NewUserService(repo repository.UserRepository, accessTokenTTL time.Duration, tokenManager auth.TokenManager, hasher hash.PasswordHasher) *UserService {
	return &UserService{
		repo:           repo,
		accessTokenTTL: accessTokenTTL,
		tokenManager:   tokenManager,
		hasher:         hasher,
	}
}

func (s *UserService) SignUp(ctx context.Context, inp SignUpUserInput) error {
	passwordHash, err := s.hasher.GeneratePasswordHash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:      inp.Name,
		Password:  passwordHash,
		Email:     inp.Email,
		CreatedAt: time.Now(),
	}

	err = s.repo.Create(ctx, user)
	return err
}

func (s *UserService) SignIn(ctx context.Context, inp SignInUserInput) (string, error) {
	user, err := s.repo.GetByEmail(ctx, inp.Email)
	if err != nil {
		return "", err
	}

	pwdValid, err := s.hasher.CheckPasswordHash(user.Password, inp.Password)
	if !pwdValid {
		return "", customErrors.ErrUserNotFound
	}
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenManager.NewJWT(user.ID, s.accessTokenTTL)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *UserService) GetByID(ctx context.Context, userID string) (domain.User, error) {
	return s.repo.GetByID(ctx, userID)
}
