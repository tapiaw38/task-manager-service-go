package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/web/middlewares"
)

func TestLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		path               string
		handler            gin.HandlerFunc
		expectedStatusCode int
	}{
		"when the handler succeeds": {
			path:               "/api/tasks",
			handler:            func(c *gin.Context) { c.Status(http.StatusOK) },
			expectedStatusCode: http.StatusOK,
		},
		"when the request carries a query string": {
			path:               "/api/tasks?completed=true",
			handler:            func(c *gin.Context) { c.Status(http.StatusOK) },
			expectedStatusCode: http.StatusOK,
		},
		"when the handler records an error": {
			path: "/api/tasks",
			handler: func(c *gin.Context) {
				_ = c.Error(assert.AnError)
				c.Status(http.StatusInternalServerError)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			router := gin.New()
			router.Use(middlewares.Logging())
			router.GET("/api/tasks", tt.handler)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tt.path, nil))

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
		})
	}
}
