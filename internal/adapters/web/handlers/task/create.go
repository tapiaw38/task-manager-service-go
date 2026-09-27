package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/web"
	"github.com/tapiaw38/task-manager-service-go/internal/usecases/task"
)

func NewCreateHandler(usecase task.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input task.CreateInput

		if err := c.ShouldBindJSON(&input); err != nil {
			web.RespondError(c, apperrors.NewApplicationError(mappings.InvalidRequestBodyError, err))

			return
		}

		output, appErr := usecase.Execute(c.Request.Context(), input)
		if appErr != nil {
			web.RespondError(c, appErr)

			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
