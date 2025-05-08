package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
	"time"
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
		if customErrors.IsDuplicateDBError(err) {
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
		if customErrors.IsDuplicateDBError(err) {
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

func (r *TaskRepo) GetListByUserID(ctx context.Context, userID string, query domain.GetTasksQuery) ([]TaskOutput, int64, error) {
	tasks := make([]TaskOutput, 0)
	var count int64

	dbQueryArgs := []any{userID}
	whereParts := make([]string, 0)
	whereArgIndex := 1

	whereParts = append(whereParts, fmt.Sprintf("t.user_id = $%d", whereArgIndex))
	whereArgIndex++

	if query.CreatedAtDateFrom != "" {
		dateFrom, err := time.Parse(time.DateOnly, query.CreatedAtDateFrom)
		if err != nil {
			return nil, 0, err
		}
		whereParts = append(whereParts, fmt.Sprintf("t.created_at >= $%d", whereArgIndex))
		dbQueryArgs = append(dbQueryArgs, dateFrom)
		whereArgIndex++
	}
	if query.CreatedAtDateTo != "" {
		dateTo, err := time.Parse(time.DateOnly, query.CreatedAtDateTo)
		if err != nil {
			return nil, 0, err
		}
		whereParts = append(whereParts, fmt.Sprintf("t.created_at <= $%d", whereArgIndex))
		dbQueryArgs = append(dbQueryArgs, dateTo)
		whereArgIndex++
	}
	if query.Completed != nil {
		whereParts = append(whereParts, fmt.Sprintf("t.completed = $%d", whereArgIndex))
		dbQueryArgs = append(dbQueryArgs, *query.Completed)
		whereArgIndex++
	}
	if len(query.CategoryIDs) > 0 {
		whereParts = append(whereParts, fmt.Sprintf("t.category_id = ANY($%d)", whereArgIndex))
		dbQueryArgs = append(dbQueryArgs, pq.Array(query.CategoryIDs))
		whereArgIndex++
	}

	whereClause := strings.Join(whereParts, " AND ")
	limitArgIndex, offsetArgIndex := whereArgIndex, whereArgIndex+1
	dbQueryArgs = append(dbQueryArgs, query.Limit, query.Offset)

	dbQuery := fmt.Sprintf(`
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
		WHERE %s
		ORDER BY t.created_at DESC
		LIMIT $%d OFFSET $%d;`, whereClause, limitArgIndex, offsetArgIndex)
	err := r.db.SelectContext(ctx, &tasks, dbQuery, dbQueryArgs...)
	if err != nil {
		return tasks, 0, err
	}

	dbQueryCountArgs := dbQueryArgs[:whereArgIndex-1] // exclude LIMIT and OFFSET
	dbQueryCount := fmt.Sprintf(`SELECT COUNT(*) FROM tasks t WHERE %s;`, whereClause)
	err = r.db.QueryRowxContext(ctx, dbQueryCount, dbQueryCountArgs...).Scan(&count)
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
