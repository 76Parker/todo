package domain

import (
	"errors"
	"strings"
	"time"
)

type Status string

var (
	statusOpen       Status = "open"
	statusInProgress Status = "in_progress"
	statusDone       Status = "done"
	statusClosed     Status = "closed"

	ErrInvalidStatus = errors.New("invalid status")
)

func NewStatus(value string) (Status, error) {
	v := strings.TrimSpace(strings.ToLower(value))
	switch Status(v) {
	case statusOpen, statusClosed, statusInProgress, statusDone:
		return Status(value), nil
	default:
		return "", ErrInvalidStatus
	}
}

type Task struct {
	id          int64
	title       string
	description string
	status      Status
	category    string
	tags        []string
	createdAt   time.Time
}

func NewTask(title, description, category string, tags []string, status Status) Task {
	normalizedTags := make([]string, len(tags))
	for i, tag := range tags {
		normalizedTags[i] = strings.TrimSpace(tag)
	}
	return Task{
		title:       normalizeString(title),
		description: normalizeString(description),
		category:    normalizeString(category),
		tags:        normalizedTags,
		status:      status,
	}
}

func (t Task) ID() int64 {
	return t.id
}
func (t Task) Title() string {
	return t.title
}
func (t Task) Description() string {
	return t.description
}
func (t Task) Category() string {
	return t.category
}
func (t Task) Tags() []string {
	return append([]string(nil), t.tags...)
}
func (t Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t Task) Status() Status {
	return t.status
}

func ReconstituteTask(id int64, title, description, status, category string, tags []string, createdAt time.Time) Task {
	return Task{
		id:          id,
		title:       title,
		description: description,
		category:    category,
		status:      Status(status),
		tags:        append([]string(nil), tags...),
		createdAt:   createdAt,
	}
}

func normalizeString(s string) string {
	return strings.TrimSpace(s)
}
