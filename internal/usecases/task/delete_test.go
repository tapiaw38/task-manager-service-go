package task_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	mock_repository "github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/repositories/task/mocks"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/task-manager-service-go/internal/usecases/task"
	"go.uber.org/mock/gomock"
)

func TestDeleteUsecase(t *testing.T) {
	type fields struct {
		repository *mock_repository.MockRepository
	}

	tests := map[string]struct {
		id           string
		prepare      func(f *fields)
		expectedErr  mappings.ErrorDetails
		expectsError bool
	}{
		"when the task is deleted": {
			id: "task-1",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "task-1").Return(nil)
			},
		},
		"when the task does not exist": {
			id: "missing",
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Delete(gomock.Any(), "missing").
					Return(apperrors.NewApplicationError(mappings.TaskNotFoundError, nil))
			},
			expectedErr:  mappings.TaskNotFoundError,
			expectsError: true,
		},
		"when the store fails": {
			id: "task-1",
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Delete(gomock.Any(), "task-1").
					Return(apperrors.NewApplicationError(mappings.TaskDeleteStoreError, errors.New("permission denied")))
			},
			expectedErr:  mappings.TaskDeleteStoreError,
			expectsError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repository: mock_repository.NewMockRepository(ctrl)}
			tt.prepare(&f)

			appErr := usecase.NewDeleteUsecase(newContextFactory(f.repository)).
				Execute(context.Background(), tt.id)

			if tt.expectsError {
				assert.NotNil(t, appErr)
				assert.True(t, appErr.IsBasedOn(tt.expectedErr))

				return
			}

			assert.Nil(t, appErr)
		})
	}
}
