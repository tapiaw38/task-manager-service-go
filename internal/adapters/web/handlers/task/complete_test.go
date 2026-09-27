package task_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handler "github.com/tapiaw38/task-manager-service-go/internal/adapters/web/handlers/task"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/task-manager-service-go/internal/usecases/task"
	mock_usecase "github.com/tapiaw38/task-manager-service-go/internal/usecases/task/mocks"
	"go.uber.org/mock/gomock"
)

func TestCompleteHandler(t *testing.T) {
	type fields struct {
		usecase *mock_usecase.MockCompleteUsecase
	}

	tests := map[string]struct {
		id                 string
		body               any
		prepare            func(f *fields)
		expectedStatusCode int
		expectedErrorCode  string
		expectedCompleted  bool
	}{
		"when the task is completed": {
			id:   "task-1",
			body: map[string]bool{"completed": true},
			prepare: func(f *fields) {
				output := newTaskOutput("task-1", "First task")
				output.Completed = true
				f.usecase.EXPECT().
					Execute(gomock.Any(), "task-1", usecase.CompleteInput{Completed: true}).
					Return(&usecase.CompleteOutput{Data: output}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedCompleted:  true,
		},
		"when the body is not valid json": {
			id:                 "task-1",
			body:               "{invalid",
			prepare:            func(f *fields) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  mappings.InvalidRequestBodyError.InternalCode,
		},
		"when the task does not exist": {
			id:   "missing",
			body: map[string]bool{"completed": true},
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), "missing", gomock.Any()).
					Return(nil, apperrors.NewApplicationError(mappings.TaskNotFoundError, nil))
			},
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  mappings.TaskNotFoundError.InternalCode,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{usecase: mock_usecase.NewMockCompleteUsecase(ctrl)}
			tt.prepare(&f)

			recorder := performRequest(t,
				http.MethodPatch,
				"/api/tasks/"+tt.id+"/complete",
				"/api/tasks/:id/complete",
				handler.NewCompleteHandler(f.usecase),
				tt.body,
			)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)

			if tt.expectedErrorCode != "" {
				assertErrorCode(t, recorder, tt.expectedErrorCode)

				return
			}

			var output usecase.CompleteOutput
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &output))
			assert.Equal(t, tt.expectedCompleted, output.Data.Completed)
		})
	}
}
