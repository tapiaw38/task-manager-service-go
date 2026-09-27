package repositories

import (
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/repositories/task"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/config"
)

type Repositories struct {
	Task task.Repository
}

type Factory func() *Repositories

func NewFactory(datasources *datasources.Datasources, configService *config.ConfigurationService) Factory {
	return func() *Repositories {
		return &Repositories{
			Task: newTaskRepository(datasources, configService),
		}
	}
}

func newTaskRepository(
	datasources *datasources.Datasources,
	configService *config.ConfigurationService,
) task.Repository {
	if configService.StoreConfig.Driver == config.FirestoreDriver {
		return task.NewFirestoreRepository(datasources.FirestoreClient)
	}

	return task.NewRepository(datasources.TaskStore)
}
