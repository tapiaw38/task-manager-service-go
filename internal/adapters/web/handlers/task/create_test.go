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

func TestCreateHandler(t *testing.T) {
	type fields struct {
		usecase *mock_usecase.MockCreateUsecase
	}

	validBody := map[string]string{
		"title":       "Comprar leche",
		"description": "Ir al supermercado",
	}

	tests := map[string]struct {
		body               any
		prepare            func(f *fields)
		expectedStatusCode int
		expectedErrorCode  string
	}{
		"when the task is created": {
			body: validBody,
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), usecase.CreateInput{
						Title:       "Comprar leche",
						Description: "Ir al supermercado",
					}).
					Return(&usecase.CreateOutput{Data: newTaskOutput("task-1", "Comprar leche")}, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		"when the body is not valid json": {
			body:               "{invalid",
			prepare:            func(f *fields) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  mappings.InvalidRequestBodyError.InternalCode,
		},
		"when the title is missing": {
			body: map[string]string{"description": "Something"},
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.NewApplicationError(mappings.TaskTitleRequiredError, nil))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  mappings.TaskTitleRequiredError.InternalCode,
		},
		"when the store fails": {
			body: validBody,
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.NewApplicationError(mappings.TaskCreateStoreError, errors.New("disk is full")))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  mappings.TaskCreateStoreError.InternalCode,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{usecase: mock_usecase.NewMockCreateUsecase(ctrl)}
			tt.prepare(&f)

			recorder := performRequest(t,
				http.MethodPost,
				"/api/tasks",
				"/api/tasks",
				handler.NewCreateHandler(f.usecase),
				tt.body,
			)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)

			if tt.expectedErrorCode != "" {
				assertErrorCode(t, recorder, tt.expectedErrorCode)

				return
			}

			var output usecase.CreateOutput
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &output))
			assert.Equal(t, "task-1", output.Data.ID)
		})
	}
}
