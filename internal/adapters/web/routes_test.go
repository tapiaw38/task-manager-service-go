package web_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/jsonstore"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/web"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/appcontext"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/config"
	"github.com/tapiaw38/task-manager-service-go/internal/usecases"
)

func newRouter(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)

	path := filepath.Join(t.TempDir(), "bd.json")
	require.NoError(t, os.WriteFile(path, []byte("[]"), 0o600))

	store, err := jsonstore.New(path)
	require.NoError(t, err)

	configService := &config.ConfigurationService{
		AppName:    "Task Manager Service",
		AppVersion: "1.0.0",
	}

	contextFactory := appcontext.NewFactory(
		datasources.CreateDatasources(store, nil),
		configService,
	)

	router := gin.New()
	web.RegisterApplicationRoutes(router, usecases.CreateUsecases(contextFactory), *configService)

	return router
}

func TestRegisterApplicationRoutes(t *testing.T) {
	tests := map[string]struct {
		method             string
		path               string
		body               string
		expectedStatusCode int
	}{
		"info is registered": {
			method:             http.MethodGet,
			path:               "/api/info",
			expectedStatusCode: http.StatusOK,
		},
		"list tasks is registered": {
			method:             http.MethodGet,
			path:               "/api/tasks",
			expectedStatusCode: http.StatusOK,
		},
		"create task is registered": {
			method:             http.MethodPost,
			path:               "/api/tasks",
			body:               `{"title":"First task","description":"Something"}`,
			expectedStatusCode: http.StatusCreated,
		},
		"get task is registered": {
			method:             http.MethodGet,
			path:               "/api/tasks/missing",
			expectedStatusCode: http.StatusNotFound,
		},
		"complete task is registered": {
			method:             http.MethodPatch,
			path:               "/api/tasks/missing/complete",
			body:               `{"completed":true}`,
			expectedStatusCode: http.StatusNotFound,
		},
		"delete task is registered": {
			method:             http.MethodDelete,
			path:               "/api/tasks/missing",
			expectedStatusCode: http.StatusNotFound,
		},
		"an unknown route answers with the not found contract": {
			method:             http.MethodGet,
			path:               "/api/unknown",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			router := newRouter(t)

			var request *http.Request
			if tt.body == "" {
				request = httptest.NewRequest(tt.method, tt.path, nil)
			} else {
				request = httptest.NewRequest(tt.method, tt.path, newBodyReader(tt.body))
				request.Header.Set("Content-Type", "application/json")
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
		})
	}
}

func newBodyReader(body string) *strings.Reader {
	return strings.NewReader(body)
}
