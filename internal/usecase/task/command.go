// Package task contain commands and queries linked with tasks
package task

// CreateCommand it's a command for service for creating new task
type CreateCommand struct {
	Title, Description, Category, Status string
	Tags                                 []string
}

// UpdateCommand it's a command for service for updating existing task
type UpdateCommand struct {
	ID                                   int64
	Title, Description, Category, Status *string
}

// TagCommand it's a command for service for adding/removing tags from task
type TagCommand struct {
	ID  int64
	Tag string
}

// QueryCommand it's a command for service for querying tasks (support only title filter)
type QueryCommand struct {
	Title string
}
