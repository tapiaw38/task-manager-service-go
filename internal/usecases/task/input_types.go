package task

import "github.com/tapiaw38/task-manager-service-go/internal/domain"

type (
	CreateInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	CompleteInput struct {
		Completed bool `json:"completed"`
	}
)

func (i CreateInput) toDomain() domain.Task {
	task := domain.Task{
		Title:       i.Title,
		Description: i.Description,
	}
	task.Sanitize()

	return task
}
