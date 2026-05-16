package task

type CreateCommand struct {
	Title, Description, Category, Status string
	Tags                                 []string
}
