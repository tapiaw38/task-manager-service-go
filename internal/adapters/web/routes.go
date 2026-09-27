package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/web/handlers/docs"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/web/handlers/info"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/web/handlers/task"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/config"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
	"github.com/tapiaw38/task-manager-service-go/internal/usecases"
)

func RegisterApplicationRoutes(
	app *gin.Engine,
	useCases *usecases.Usecases,
	configService config.ConfigurationService,
) {
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, apperrors.NewApplicationError(mappings.NotFoundError, nil))
	})

	api := app.Group("/api")

	api.GET("/info", info.NewGetHandler(configService))
	api.GET("/docs", docs.NewUIHandler())
	api.GET("/docs/openapi.yaml", docs.NewSpecHandler())

	tasks := api.Group("/tasks")
	tasks.GET("", task.NewListHandler(useCases.Task.ListUsecase))
	tasks.POST("", task.NewCreateHandler(useCases.Task.CreateUsecase))
	tasks.GET("/:id", task.NewGetHandler(useCases.Task.GetUsecase))
	tasks.PATCH("/:id/complete", task.NewCompleteHandler(useCases.Task.CompleteUsecase))
	tasks.DELETE("/:id", task.NewDeleteHandler(useCases.Task.DeleteUsecase))
}
