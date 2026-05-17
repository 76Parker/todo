package task

import (
	"context"
	"todo/internal/entities/domain"
)

// Service handle commands and queries, create/read domain.Task and interact with persistence storage (Repository)
type Service struct {
	repo Repository
}

// NewService it's a constructor for Service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create handle CreateCommand use case
func (t *Service) Create(ctx context.Context, cmd CreateCommand) (domain.Task, error) {

	status, err := domain.NewStatus(cmd.Status)
	if err != nil {
		return domain.Task{}, err
	}
	task := domain.NewTask(cmd.Title, cmd.Description, cmd.Category, cmd.Tags, status)
	return t.repo.Create(ctx, task)
}

// ReadByID read task from storage by id
func (t *Service) ReadByID(ctx context.Context, id int64) (domain.Task, error) {
	return t.repo.ReadByID(ctx, id)
}
