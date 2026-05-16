package task

import (
	"context"
	"todo/internal/entities/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (t *Service) Create(ctx context.Context, cmd CreateCommand) (domain.Task, error) {

	status, err := domain.NewStatus(cmd.Status)
	if err != nil {
		return domain.Task{}, err
	}
	task := domain.NewTask(cmd.Title, cmd.Description, cmd.Category, cmd.Tags, status)
	return t.repo.Create(ctx, task)
}
