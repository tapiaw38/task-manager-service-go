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

func TestGetHandler(t *testing.T) {
	type fields struct {
		usecase *mock_usecase.MockGetUsecase
	}

	tests := map[string]struct {
		id                 string
		prepare            func(f *fields)
		expectedStatusCode int
		expectedErrorCode  string
	}{
		"when the task exists": {
			id: "task-1",
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), "task-1").
					Return(&usecase.GetOutput{Data: newTaskOutput("task-1", "First task")}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		"when the task does not exist": {
			id: "missing",
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), "missing").
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

			f := fields{usecase: mock_usecase.NewMockGetUsecase(ctrl)}
			tt.prepare(&f)

			recorder := performRequest(t,
				http.MethodGet,
				"/api/tasks/"+tt.id,
				"/api/tasks/:id",
				handler.NewGetHandler(f.usecase),
				nil,
			)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)

			if tt.expectedErrorCode != "" {
				assertErrorCode(t, recorder, tt.expectedErrorCode)

				return
			}

			var output usecase.GetOutput
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &output))
			assert.Equal(t, "task-1", output.Data.ID)
		})
	}
}
