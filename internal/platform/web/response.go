package web

import (
	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
)

func RespondError(c *gin.Context, appErr apperrors.ApplicationError) {
	appErr.Log(c.Request.Context())
	c.JSON(appErr.StatusCode(), appErr)
}
