APP_NAME=todo
MAIN=./cmd

.PHONY: test


docker-build:
	docker build -t $(APP_NAME) .

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

run:
	golangci-lint run
	go run $(MAIN)

test:
	go test -v ./...

mocks:
	mockgen -source=internal/api/task_handler.go -destination=internal/api/task_handler_mock_test.go -package api
	mockgen -source=internal/usecase/task/repository.go -destination=internal/usecase/task/repository_mock_test.go -package task

build:
	GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) $(MAIN)

arm-build:
	GOOS=darwin GOARCH=arm64 go build -o $(APP_NAME) $(MAIN)
