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

// UpdateByID handle UpdateCommand use case
func (t *Service) UpdateByID(ctx context.Context, cmd UpdateCommand) (domain.Task, error) {
	if cmd.Status != nil {
		status, err := domain.NewStatus(*cmd.Status)
		if err != nil {
			return domain.Task{}, err
		}
		s := string(status)
		cmd.Status = &s
	}
	return t.repo.UpdateByID(ctx, cmd)
}

// AddTagByID handle TagCommand use case for adding tag to task
func (t *Service) AddTagByID(ctx context.Context, cmd TagCommand) (domain.Task, error) {
	return t.repo.AddTagByID(ctx, cmd)
}

// DeleteTagByID handle TagCommand use case for removing tag from task
func (t *Service) DeleteTagByID(ctx context.Context, cmd TagCommand) (domain.Task, error) {
	return t.repo.DeleteTagByID(ctx, cmd)
}

// QueryTasks search tasks by title
func (t *Service) QueryTasks(ctx context.Context, cmd QueryCommand) ([]domain.Task, error) {
	return t.repo.QueryTasks(ctx, cmd)
}

// DeleteTaskByID delete task from storage by id
func (t *Service) DeleteTaskByID(ctx context.Context, id int64) error {
	return t.repo.DeleteByID(ctx, id)
}
