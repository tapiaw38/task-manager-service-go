package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/web"
	"github.com/tapiaw38/task-manager-service-go/internal/usecases/task"
)

func NewDeleteHandler(usecase task.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := usecase.Execute(c.Request.Context(), c.Param("id")); appErr != nil {
			web.RespondError(c, appErr)

			return
		}

		c.Status(http.StatusNoContent)
	}
}
