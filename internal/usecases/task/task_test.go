package task_test

import (
	"time"

	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/repositories"
	mock_repository "github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/repositories/task/mocks"
	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/appcontext"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/config"
)

var referenceTime = time.Date(2026, time.September, 27, 10, 5, 0, 0, time.UTC)

func newContextFactory(repo *mock_repository.MockRepository) appcontext.Factory {
	return func(opts ...appcontext.Option) *appcontext.Context {
		return &appcontext.Context{
			Repositories:  &repositories.Repositories{Task: repo},
			ConfigService: &config.ConfigurationService{},
		}
	}
}

func newTask(id, title string) domain.Task {
	return domain.Task{
		ID:          id,
		Title:       title,
		Description: "Description of " + title,
		Completed:   false,
		CreatedAt:   referenceTime,
		UpdatedAt:   referenceTime,
	}
}
