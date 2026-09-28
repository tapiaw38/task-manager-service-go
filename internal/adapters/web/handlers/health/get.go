package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Output struct {
	Status string `json:"status"`
}

func NewGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, Output{Status: "ok"})
	}
}
