// Package domain provide domain struct and validating constructors for domain structs
package domain

import (
	"errors"
	"strings"
	"time"
)

// Status domain status for Task
type Status string

var (
	statusOpen       Status = "open"
	statusInProgress Status = "in_progress"
	statusDone       Status = "done"
	statusClosed     Status = "closed"

	// ErrInvalidStatus error for invalid statuses
	ErrInvalidStatus = errors.New("invalid status")
	// ErrTaskNotFound error for situation where task not found in storage
	ErrTaskNotFound = errors.New("task not found")

	// ErrNoChanges error for no changes in task update
	ErrNoChanges = errors.New("no changes in task update")

	// ErrTagOrTaskNotFound error for situation where tag or task not found in storage
	ErrTagOrTaskNotFound = errors.New("tag or task not found")
)

// NewStatus construct Status and check domain rules (check valid statuses for Task)
func NewStatus(value string) (Status, error) {
	v := strings.TrimSpace(strings.ToLower(value))
	switch Status(v) {
	case statusOpen, statusClosed, statusInProgress, statusDone:
		return Status(value), nil
	default:
		return "", ErrInvalidStatus
	}
}

// Task it's a domain TODO-task
type Task struct {
	id          int64
	title       string
	description string
	status      Status
	category    string
	tags        []string
	createdAt   time.Time
}

// NewTask constructor for Task
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

// ID return task id
func (t Task) ID() int64 {
	return t.id
}

// Title return task title
func (t Task) Title() string {
	return t.title
}

// Description return task description
func (t Task) Description() string {
	return t.description
}

// Category return task category
func (t Task) Category() string {
	return t.category
}

// Tags return task tags
func (t Task) Tags() []string {
	return append([]string(nil), t.tags...)
}

// CreatedAt return task creation time
func (t Task) CreatedAt() time.Time {
	return t.createdAt
}

// Status return task status (open, in_progress, done, closed)
func (t Task) Status() Status {
	return t.status
}

// ReconstituteTask method for repositories for restoring task from DB
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
