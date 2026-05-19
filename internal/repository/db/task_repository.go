// Package db provide implementations for repositories
package db

import (
	"context"
	"errors"
	"time"
	"todo/internal/entities/domain"
	"todo/internal/entities/table"
	"todo/internal/usecase/task"

	"github.com/Masterminds/squirrel"
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

// ReadByID extract task from DB
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

// UpdateByID updates task by ID and returns updated task
func (t *TaskRepository) UpdateByID(ctx context.Context, updateCmd task.UpdateCommand) (domain.Task, error) {
	qb := squirrel.Update("todo.tasks").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": updateCmd.ID})

	qb, hasChanges := buildUpdateQuery(qb, updateCmd)
	if !hasChanges {
		return domain.Task{}, domain.ErrNoChanges
	}
	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		return domain.Task{}, err
	}
	rows, err := t.pool.Query(ctx, sqlQuery, args...)
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
	domainTask := domain.ReconstituteTask(
		task.ID,
		task.Title,
		task.Description.String,
		task.Status,
		task.Category.String,
		pgTextArrayToStrings(task.Tags),
		task.CreatedAt)
	return domainTask, nil
}

func buildUpdateQuery(qb squirrel.UpdateBuilder, updateCmd task.UpdateCommand) (squirrel.UpdateBuilder, bool) {
	var hasChanges bool
	if updateCmd.Title != nil {
		qb = qb.Set("title", *updateCmd.Title)
		hasChanges = true
	}
	if updateCmd.Description != nil {
		qb = qb.Set("description", *updateCmd.Description)
		hasChanges = true
	}
	if updateCmd.Category != nil {
		qb = qb.Set("category", *updateCmd.Category)
		hasChanges = true
	}
	if updateCmd.Status != nil {
		qb = qb.Set("status", *updateCmd.Status)
		hasChanges = true
	}
	return qb.Suffix("RETURNING id, title, status, description, category, tags, created_at"), hasChanges
}

// AddTagByID adds a tag to a task by task ID
func (t *TaskRepository) AddTagByID(ctx context.Context, cmd task.TagCommand) (domain.Task, error) {

	sqlQuery := `UPDATE todo.tasks SET tags = array_append(tags, $1) WHERE id = $2 RETURNING id, title, status, description, category, tags, created_at`
	rows, err := t.pool.Query(ctx, sqlQuery, cmd.Tag, cmd.ID)
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
	domainTask := domain.ReconstituteTask(
		task.ID,
		task.Title,
		task.Description.String,
		task.Status,
		task.Category.String,
		pgTextArrayToStrings(task.Tags),
		task.CreatedAt)
	return domainTask, nil
}

// DeleteTagByID deletes tags from a task by task ID
func (t *TaskRepository) DeleteTagByID(ctx context.Context, cmd task.TagCommand) (domain.Task, error) {
	sqlQuery := `UPDATE todo.tasks SET tags = array_remove(tags, $1)
				WHERE id = $2 AND $1 = ANY(tags)
				RETURNING id, title, status, description, category, tags, created_at`
	rows, err := t.pool.Query(ctx, sqlQuery, cmd.Tag, cmd.ID)
	if err != nil {
		return domain.Task{}, err
	}
	task, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[table.Task])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrTagOrTaskNotFound
		}
		return domain.Task{}, err
	}
	domainTask := domain.ReconstituteTask(
		task.ID,
		task.Title,
		task.Description.String,
		task.Status,
		task.Category.String,
		pgTextArrayToStrings(task.Tags),
		task.CreatedAt)
	return domainTask, nil
}

// QueryTasks returns a list of tasks that match the given query command.
func (t *TaskRepository) QueryTasks(ctx context.Context, cmd task.QueryCommand) ([]domain.Task, error) {
	sqlQuery := `SELECT id, title, status, description, category, tags, created_at
				FROM todo.tasks`

	args := make([]any, 0)
	if cmd.Title != "" {
		sqlQuery += ` WHERE title ILIKE $1`
		args = append(args, "%"+cmd.Title+"%")
	}
	rows, err := t.pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[table.Task])
	if err != nil {
		return nil, err
	}
	domainTasks := make([]domain.Task, 0, len(tasks))
	for _, task := range tasks {
		domainTasks = append(domainTasks, domain.ReconstituteTask(
			task.ID,
			task.Title,
			task.Description.String,
			task.Status,
			task.Category.String,
			pgTextArrayToStrings(task.Tags),
			task.CreatedAt))
	}
	return domainTasks, nil
}

// DeleteByID deletes a task by its ID.
func (t *TaskRepository) DeleteByID(ctx context.Context, id int64) error {
	sqlQuery := `DELETE FROM todo.tasks WHERE id = $1`
	cmd, err := t.pool.Exec(ctx, sqlQuery, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}
