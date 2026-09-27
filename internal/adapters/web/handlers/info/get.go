package info

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/config"
)

type Output struct {
	Application string `json:"application"`
	Version     string `json:"version"`
}

func NewGetHandler(configService config.ConfigurationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, Output{
			Application: configService.AppName,
			Version:     configService.AppVersion,
		})
	}
}
