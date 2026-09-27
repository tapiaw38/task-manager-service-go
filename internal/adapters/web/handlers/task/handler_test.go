package task_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	usecase "github.com/tapiaw38/task-manager-service-go/internal/usecases/task"
)

var referenceTime = time.Date(2026, time.September, 27, 10, 5, 0, 0, time.UTC)

func init() {
	gin.SetMode(gin.TestMode)
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func performRequest(t *testing.T, method, path, route string, handler gin.HandlerFunc, body any) *httptest.ResponseRecorder {
	t.Helper()

	router := gin.New()
	router.Handle(method, route, handler)

	var payload *bytes.Reader

	switch typed := body.(type) {
	case nil:
		payload = bytes.NewReader(nil)
	case string:
		payload = bytes.NewReader([]byte(typed))
	default:
		encoded, err := json.Marshal(typed)
		require.NoError(t, err)
		payload = bytes.NewReader(encoded)
	}

	request := httptest.NewRequest(method, path, payload)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}

func assertErrorCode(t *testing.T, recorder *httptest.ResponseRecorder, expectedCode string) {
	t.Helper()

	var response errorResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, expectedCode, response.Code)
	assert.NotEmpty(t, response.Message)
}

func newTaskOutput(id, title string) usecase.TaskOutputData {
	return usecase.TaskOutputData{
		ID:          id,
		Title:       title,
		Description: "Description of " + title,
		Completed:   false,
		CreatedAt:   referenceTime,
		UpdatedAt:   referenceTime,
	}
}
