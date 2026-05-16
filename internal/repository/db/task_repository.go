package db

import (
	"context"
	"time"
	"todo/internal/entities/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

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
