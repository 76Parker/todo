// Package task contain commands and queries linked with TODO-tasks
package task

// CreateCommand it's a command for service for creating new task
type CreateCommand struct {
	Title, Description, Category, Status string
	Tags                                 []string
}
