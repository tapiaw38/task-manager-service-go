package errors_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

func TestApplicationErrorFields(t *testing.T) {
	tests := map[string]struct {
		details         mappings.ErrorDetails
		cause           error
		expectedCode    string
		expectedStatus  int
		expectedMessage string
		expectedOrigin  string
	}{
		"when the error carries an original cause": {
			details:         mappings.TaskNotFoundError,
			cause:           errors.New("document missing"),
			expectedCode:    "task:shared:not-found",
			expectedStatus:  http.StatusNotFound,
			expectedMessage: "task not found",
			expectedOrigin:  "document missing",
		},
		"when the error has no original cause": {
			details:         mappings.InvalidRequestBodyError,
			cause:           nil,
			expectedCode:    "common:request-body-parsing-error",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "invalid request body format",
			expectedOrigin:  "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			appErr := apperrors.NewApplicationError(tt.details, tt.cause)

			assert.Equal(t, tt.expectedCode, appErr.InternalCode())
			assert.Equal(t, tt.expectedStatus, appErr.StatusCode())
			assert.Equal(t, tt.expectedMessage, appErr.Message())
			assert.Equal(t, tt.expectedOrigin, appErr.OriginalMessage())
			assert.Equal(t, tt.expectedMessage, appErr.Error())
		})
	}
}

func TestApplicationErrorIsBasedOn(t *testing.T) {
	appErr := apperrors.NewApplicationError(mappings.TaskNotFoundError, nil)

	assert.True(t, appErr.IsBasedOn(mappings.TaskNotFoundError))
	assert.False(t, appErr.IsBasedOn(mappings.TaskCreateStoreError))
}

func TestApplicationErrorMarshalJSON(t *testing.T) {
	appErr := apperrors.NewApplicationError(mappings.TaskTitleRequiredError, errors.New("hidden cause"))

	encoded, err := json.Marshal(appErr)
	require.NoError(t, err)

	var body map[string]string
	require.NoError(t, json.Unmarshal(encoded, &body))

	assert.Equal(t, map[string]string{
		"code":    "task:validation:title-required",
		"message": "task title is required",
	}, body)
}

func TestApplicationErrorAddExtraFields(t *testing.T) {
	appErr := apperrors.NewApplicationError(mappings.TaskNotFoundError, nil).
		AddExtraFields(map[string]any{"task_id": "task-1"}).
		AddExtraFields(map[string]any{"attempt": 2})

	assert.True(t, appErr.IsBasedOn(mappings.TaskNotFoundError))

	encoded, err := json.Marshal(appErr)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), "task_id")
}

func TestApplicationErrorLog(t *testing.T) {
	appErr := apperrors.NewApplicationError(mappings.TaskDeleteStoreError, errors.New("permission denied")).
		AddExtraFields(map[string]any{"task_id": "task-1"})

	assert.NotPanics(t, func() {
		appErr.Log(context.Background())
	})
}
