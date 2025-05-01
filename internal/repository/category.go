package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"todo_list_go/internal/domain"
	customErrors "todo_list_go/pkg/errors"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	var createdCategory domain.Category
	query := `
		INSERT INTO categories (user_id, created_at, title, description, color) 
		values ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, title, description, color;`
	err := r.db.QueryRowxContext(
		ctx, query, category.UserID, category.CreatedAt, category.Title, category.Description, category.Color,
	).StructScan(&createdCategory)

	if err != nil {
		if customErrors.IsDuplicateKeyError(err) {
			return domain.Category{}, customErrors.ErrCategoryAlreadyExists
		}
		return domain.Category{}, err
	}
	return createdCategory, nil
}

func (r *CategoryRepo) Update(ctx context.Context, inp UpdateCategoryInput) (domain.Category, error) {
	var updatedCategory domain.Category

	setClause := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if inp.Title != nil {
		setClause = append(setClause, fmt.Sprintf("title = $%d", argID))
		args = append(args, inp.Title)
		argID++
	}
	if inp.Description != nil {
		setClause = append(setClause, fmt.Sprintf("description = $%d", argID))
		args = append(args, inp.Description)
		argID++
	}
	if inp.Color != nil {
		setClause = append(setClause, fmt.Sprintf("color = $%d", argID))
		args = append(args, inp.Color)
		argID++
	}

	setQuery := strings.Join(setClause, ", ")
	query := fmt.Sprintf(
		`UPDATE categories SET %s WHERE id = $%d 
                RETURNING id, user_id, created_at, title, description, color;`,
		setQuery, argID,
	)
	args = append(args, inp.ID)

	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&updatedCategory)
	if err != nil {
		if customErrors.IsDuplicateKeyError(err) {
			return domain.Category{}, customErrors.ErrCategoryAlreadyExists
		}
		return domain.Category{}, err
	}

	return updatedCategory, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM categories WHERE id=$1;"
	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *CategoryRepo) GetListByUserID(ctx context.Context, userID string) ([]domain.Category, error) {
	categories := make([]domain.Category, 0)

	query := "SELECT id, created_at, title, description, color FROM categories WHERE user_id=$1 ORDER BY created_at DESC;"
	err := r.db.SelectContext(ctx, &categories, query, userID)

	return categories, err
}

func (r *CategoryRepo) GetByID(ctx context.Context, categoryID, userID string) (domain.Category, error) {
	var category domain.Category

	query := "SELECT id, created_at, title, description, color FROM categories WHERE id=$1 AND user_id=$2;"
	err := r.db.GetContext(ctx, &category, query, categoryID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Category{}, customErrors.ErrCategoryNotFound
		}

		return domain.Category{}, err
	}

	return category, nil
}
