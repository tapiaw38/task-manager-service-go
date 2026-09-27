package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/web"
	"github.com/tapiaw38/task-manager-service-go/internal/usecases/task"
)

func NewGetHandler(usecase task.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := usecase.Execute(c.Request.Context(), c.Param("id"))
		if appErr != nil {
			web.RespondError(c, appErr)

			return
		}

		c.JSON(http.StatusOK, output)
	}
}
