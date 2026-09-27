package task_test

import (
	"encoding/json"
	"errors"
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

func TestListHandler(t *testing.T) {
	type fields struct {
		usecase *mock_usecase.MockListUsecase
	}

	tests := map[string]struct {
		prepare            func(f *fields)
		expectedStatusCode int
		expectedErrorCode  string
		expectedTotal      int
	}{
		"when there are tasks": {
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any()).
					Return(&usecase.ListOutput{
						Data:  []usecase.TaskOutputData{newTaskOutput("task-1", "First task")},
						Total: 1,
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedTotal:      1,
		},
		"when there are no tasks": {
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any()).
					Return(&usecase.ListOutput{Data: []usecase.TaskOutputData{}, Total: 0}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedTotal:      0,
		},
		"when the store fails": {
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any()).
					Return(nil, apperrors.NewApplicationError(mappings.TaskListStoreError, errors.New("corrupted file")))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  mappings.TaskListStoreError.InternalCode,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{usecase: mock_usecase.NewMockListUsecase(ctrl)}
			tt.prepare(&f)

			recorder := performRequest(t,
				http.MethodGet,
				"/api/tasks",
				"/api/tasks",
				handler.NewListHandler(f.usecase),
				nil,
			)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)

			if tt.expectedErrorCode != "" {
				assertErrorCode(t, recorder, tt.expectedErrorCode)

				return
			}

			var output usecase.ListOutput
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &output))
			assert.Equal(t, tt.expectedTotal, output.Total)
		})
	}
}
