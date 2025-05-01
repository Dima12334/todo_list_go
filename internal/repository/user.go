package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"todo_list_go/internal/domain"
	customErrors "todo_list_go/pkg/errors"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) error {
	query := "INSERT INTO users (name, email, password, created_at) values ($1, $2, $3, $4);"
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		if customErrors.IsDuplicateKeyError(err) {
			return customErrors.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	query := "SELECT id, created_at, name, email FROM users WHERE id = $1;"
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, customErrors.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	query := "SELECT id, created_at, name, email, password FROM users WHERE email = $1;"
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, customErrors.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}
