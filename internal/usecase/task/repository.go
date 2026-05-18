package task

import (
	"context"
	"todo/internal/entities/domain"
)

// Repository it's interface to interact with persistence storage
type Repository interface {
	Create(ctx context.Context, task domain.Task) (domain.Task, error)
	ReadByID(ctx context.Context, id int64) (domain.Task, error)
	UpdateByID(ctx context.Context, updateCmd UpdateCommand) (domain.Task, error)
}
