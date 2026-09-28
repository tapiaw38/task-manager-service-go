package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handler "github.com/tapiaw38/task-manager-service-go/internal/adapters/web/handlers/health"
)

func TestGetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/health", handler.NewGetHandler())

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))

	assert.Equal(t, http.StatusOK, recorder.Code)

	var output handler.Output
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &output))
	assert.Equal(t, "ok", output.Status)
}
