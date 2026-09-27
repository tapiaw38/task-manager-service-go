package appcontext

import (
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/config"
)

type Context struct {
	Repositories  *repositories.Repositories
	ConfigService *config.ConfigurationService
}

type Option func(*Context)

type Factory func(opts ...Option) *Context

func NewFactory(
	datasources *datasources.Datasources,
	configService *config.ConfigurationService,
) Factory {
	repositoriesFactory := repositories.NewFactory(datasources, configService)

	return func(opts ...Option) *Context {
		ctx := &Context{
			Repositories:  repositoriesFactory(),
			ConfigService: configService,
		}

		for _, opt := range opts {
			opt(ctx)
		}

		return ctx
	}
}
