// Package db provide implementations for repositories
package db

import (
	"context"
	"errors"
	"time"
	"todo/internal/entities/domain"
	"todo/internal/entities/table"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TaskRepository provide CRUD-methods for domain.Task
type TaskRepository struct {
	pool *pgxpool.Pool
}

// NewTaskRepository it's a constructor for TaskRepository
func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

// Create new domain.Task in database
func (t *TaskRepository) Create(ctx context.Context, task domain.Task) (domain.Task, error) {

	sqlQuery := `INSERT INTO todo.tasks 
    (title, status, description,category, tags) 
	 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	var id int64
	var createdAt time.Time
	if err := t.pool.QueryRow(ctx, sqlQuery,
		task.Title(),
		string(task.Status()),
		task.Description(),
		task.Category(),
		task.Tags()).Scan(&id, &createdAt); err != nil {
		return domain.Task{}, err
	}
	restoredTask := domain.ReconstituteTask(
		id,
		task.Title(),
		task.Description(),
		string(task.Status()),
		task.Category(),
		task.Tags(),
		createdAt,
	)
	return restoredTask, nil
}

// ReadByID extract TODO-task from DB
func (t *TaskRepository) ReadByID(ctx context.Context, id int64) (domain.Task, error) {
	sqlQuery := `SELECT id, title, description, status, category, tags, created_at 
					 FROM todo.tasks WHERE id = $1`

	rows, err := t.pool.Query(ctx, sqlQuery, id)
	if err != nil {
		return domain.Task{}, err
	}

	task, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[table.Task])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}
		return domain.Task{}, err
	}
	restoredTask := domain.ReconstituteTask(
		task.ID,
		task.Title,
		task.Description.String,
		task.Status,
		task.Category.String,
		pgTextArrayToStrings(task.Tags),
		task.CreatedAt,
	)
	return restoredTask, nil
}

func pgTextArrayToStrings(arr pgtype.Array[pgtype.Text]) []string {

	result := make([]string, 0, len(arr.Elements))
	for _, elem := range arr.Elements {
		if elem.Valid {
			result = append(result, elem.String)
		}
	}
	return result
}
