package task_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	handler "github.com/tapiaw38/task-manager-service-go/internal/adapters/web/handlers/task"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
	mock_usecase "github.com/tapiaw38/task-manager-service-go/internal/usecases/task/mocks"
	"go.uber.org/mock/gomock"
)

func TestDeleteHandler(t *testing.T) {
	type fields struct {
		usecase *mock_usecase.MockDeleteUsecase
	}

	tests := map[string]struct {
		id                 string
		prepare            func(f *fields)
		expectedStatusCode int
		expectedErrorCode  string
	}{
		"when the task is deleted": {
			id: "task-1",
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), "task-1").Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		"when the task does not exist": {
			id: "missing",
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), "missing").
					Return(apperrors.NewApplicationError(mappings.TaskNotFoundError, nil))
			},
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  mappings.TaskNotFoundError.InternalCode,
		},
		"when the store fails": {
			id: "task-1",
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), "task-1").
					Return(apperrors.NewApplicationError(mappings.TaskDeleteStoreError, errors.New("permission denied")))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  mappings.TaskDeleteStoreError.InternalCode,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{usecase: mock_usecase.NewMockDeleteUsecase(ctrl)}
			tt.prepare(&f)

			recorder := performRequest(t,
				http.MethodDelete,
				"/api/tasks/"+tt.id,
				"/api/tasks/:id",
				handler.NewDeleteHandler(f.usecase),
				nil,
			)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)

			if tt.expectedErrorCode != "" {
				assertErrorCode(t, recorder, tt.expectedErrorCode)

				return
			}

			assert.Empty(t, recorder.Body.String())
		})
	}
}
