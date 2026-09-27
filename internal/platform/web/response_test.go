package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/web"
)

func TestRespondError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		details            mappings.ErrorDetails
		expectedStatusCode int
		expectedCode       string
	}{
		"when the task is not found": {
			details:            mappings.TaskNotFoundError,
			expectedStatusCode: http.StatusNotFound,
			expectedCode:       "task:shared:not-found",
		},
		"when the body could not be parsed": {
			details:            mappings.InvalidRequestBodyError,
			expectedStatusCode: http.StatusBadRequest,
			expectedCode:       "common:request-body-parsing-error",
		},
		"when the store fails": {
			details:            mappings.TaskListStoreError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedCode:       "task:list:store-error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			router := gin.New()
			router.GET("/error", func(c *gin.Context) {
				web.RespondError(c, apperrors.NewApplicationError(tt.details, nil))
			})

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/error", nil))

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)

			var body map[string]string
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
			assert.Equal(t, tt.expectedCode, body["code"])
			assert.NotEmpty(t, body["message"])
		})
	}
}
