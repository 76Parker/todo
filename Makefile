coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test:
	go test -v ./...

mocks:
	mockgen -source=internal/api/task_handler.go -destination=internal/api/task_handler_mock.go -package api
	mockgen -source=internal/usecase/task/repository.go -destination=internal/usecase/task/repository_mock.go -package task
