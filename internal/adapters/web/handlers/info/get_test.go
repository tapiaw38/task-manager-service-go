package info_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handler "github.com/tapiaw38/task-manager-service-go/internal/adapters/web/handlers/info"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/config"
)

func TestGetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/info", handler.NewGetHandler(config.ConfigurationService{
		AppName:    "Task Manager Service",
		AppVersion: "1.0.0",
	}))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/info", nil))

	assert.Equal(t, http.StatusOK, recorder.Code)

	var output handler.Output
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &output))

	assert.Equal(t, "Task Manager Service", output.Application)
	assert.Equal(t, "1.0.0", output.Version)
}
