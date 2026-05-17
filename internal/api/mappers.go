// Package api provide API-handlers for interact with the application
package api

import (
	"todo/internal/entities/domain"
	"todo/internal/entities/dto"
	"todo/internal/usecase/task"
)

func createTaskDtoToCommand(taskDto dto.CreateTask) task.CreateCommand {
	var description, category string
	var tags []string

	if taskDto.Description != nil {
		description = *taskDto.Description
	}
	if taskDto.Category != nil {
		category = *taskDto.Category
	}
	if taskDto.Tags != nil {
		tags = append([]string(nil), *taskDto.Tags...)
	}
	return task.CreateCommand{
		Title:       taskDto.Title,
		Description: description,
		Status:      taskDto.Status,
		Category:    category,
		Tags:        tags,
	}
}

func fromDomainTaskToReadTaskDTO(task domain.Task) dto.ReadTask {
	return dto.ReadTask{
		ID:          task.ID(),
		Title:       task.Title(),
		Status:      string(task.Status()),
		Description: task.Description(),
		Category:    task.Category(),
		Tags:        task.Tags(),
		CreatedAt:   task.CreatedAt(),
	}
}
