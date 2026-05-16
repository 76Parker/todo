package task

import (
	"context"
	"todo/internal/entities/domain"
)

type Repository interface {
	Create(ctx context.Context, task domain.Task) (domain.Task, error)
}
