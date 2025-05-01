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

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, task domain.Task) (TaskOutput, error) {
	var createdTaskID string
	var createdTask TaskOutput

	query := `
		INSERT INTO tasks (created_at, updated_at, user_id, category_id, title, description, completed)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`
	err := r.db.QueryRowxContext(
		ctx, query, task.CreatedAt, task.UpdatedAt, task.UserID, task.CategoryID, task.Title, task.Description, task.Completed,
	).Scan(&createdTaskID)
	if err != nil {
		if customErrors.IsDuplicateKeyError(err) {
			return TaskOutput{}, customErrors.ErrTaskAlreadyExists
		}
		return TaskOutput{}, err
	}

	query = `
		SELECT     
		t.id AS id,
		t.created_at AS created_at,
		t.updated_at AS updated_at,
		t.title AS title,
		t.description AS description,
		t.completed AS completed,
		
		c.id AS "category.id",
		c.created_at AS "category.created_at",
		c.title AS "category.title",
		c.description AS "category.description",
		c.color AS "category.color"
		FROM tasks t
		INNER JOIN categories c ON t.category_id = c.id 
		WHERE t.id = $1;`
	err = r.db.QueryRowxContext(ctx, query, createdTaskID).StructScan(&createdTask)
	if err != nil {
		return TaskOutput{}, err
	}

	return createdTask, nil
}

func (r *TaskRepo) Update(ctx context.Context, inp UpdateTaskInput) (TaskOutput, error) {
	var updatedTaskID string

	setClause := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	setClause = append(setClause, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, inp.UpdatedAt)
	argID++
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
	if inp.CategoryID != nil {
		setClause = append(setClause, fmt.Sprintf("category_id = $%d", argID))
		args = append(args, inp.CategoryID)
		argID++
	}
	if inp.Completed != nil {
		setClause = append(setClause, fmt.Sprintf("completed = $%d", argID))
		args = append(args, inp.Completed)
		argID++
	}

	setQuery := strings.Join(setClause, ", ")
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d RETURNING id;", setQuery, argID)
	args = append(args, inp.ID)
	err := r.db.QueryRowxContext(ctx, query, args...).Scan(&updatedTaskID)
	if err != nil {
		if customErrors.IsDuplicateKeyError(err) {
			return TaskOutput{}, customErrors.ErrTaskAlreadyExists
		}
		return TaskOutput{}, err
	}

	return r.GetByID(ctx, updatedTaskID, inp.UserID)
}

func (r *TaskRepo) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM tasks WHERE id = $1;"
	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *TaskRepo) GetListByUserID(ctx context.Context, userID string, limit, offset int) ([]TaskOutput, int64, error) {
	tasks := make([]TaskOutput, 0)
	var count int64

	query := `
		SELECT     
		t.id AS id,
		t.created_at AS created_at,
		t.updated_at AS updated_at,
		t.title AS title,
		t.description AS description,
		t.completed AS completed,
		
		c.id AS "category.id",
		c.created_at AS "category.created_at",
		c.title AS "category.title",
		c.description AS "category.description",
		c.color AS "category.color"
		FROM tasks t
		INNER JOIN categories c ON t.category_id = c.id 
		WHERE t.user_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3;`
	err := r.db.SelectContext(ctx, &tasks, query, userID, limit, offset)
	if err != nil {
		return tasks, 0, err
	}

	queryCount := `SELECT COUNT(*) FROM tasks WHERE user_id = $1;`
	err = r.db.QueryRowxContext(ctx, queryCount, userID).Scan(&count)
	if err != nil {
		return tasks, 0, err
	}

	return tasks, count, err
}

func (r *TaskRepo) GetByID(ctx context.Context, taskID, userID string) (TaskOutput, error) {
	var task TaskOutput

	query := `
		SELECT     
		t.id AS id,
		t.created_at AS created_at,
		t.updated_at AS updated_at,
		t.title AS title,
		t.description AS description,
		t.completed AS completed,
		
		c.id AS "category.id",
		c.created_at AS "category.created_at",
		c.title AS "category.title",
		c.description AS "category.description",
		c.color AS "category.color"
		FROM tasks t
		INNER JOIN categories c ON t.category_id = c.id 
		WHERE t.id = $1 AND t.user_id = $2;`

	err := r.db.QueryRowxContext(ctx, query, taskID, userID).StructScan(&task)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TaskOutput{}, customErrors.ErrTaskNotFound
		}

		return TaskOutput{}, err
	}

	return task, nil
}
