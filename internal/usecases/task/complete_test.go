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

func TestCompleteUsecase(t *testing.T) {
	type fields struct {
		repository *mock_repository.MockRepository
	}

	tests := map[string]struct {
		id                string
		input             usecase.CompleteInput
		prepare           func(f *fields)
		expectedCompleted bool
		expectedErr       mappings.ErrorDetails
		expectsError      bool
	}{
		"when the task is marked as completed": {
			id:    "task-1",
			input: usecase.CompleteInput{Completed: true},
			prepare: func(f *fields) {
				completed := newTask("task-1", "First task")
				completed.Completed = true
				f.repository.EXPECT().
					SetCompleted(gomock.Any(), "task-1", true).
					Return(&completed, nil)
			},
			expectedCompleted: true,
		},
		"when the task is marked as pending again": {
			id:    "task-1",
			input: usecase.CompleteInput{Completed: false},
			prepare: func(f *fields) {
				pending := newTask("task-1", "First task")
				f.repository.EXPECT().
					SetCompleted(gomock.Any(), "task-1", false).
					Return(&pending, nil)
			},
			expectedCompleted: false,
		},
		"when the task does not exist": {
			id:    "missing",
			input: usecase.CompleteInput{Completed: true},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					SetCompleted(gomock.Any(), "missing", true).
					Return(nil, apperrors.NewApplicationError(mappings.TaskNotFoundError, nil))
			},
			expectedErr:  mappings.TaskNotFoundError,
			expectsError: true,
		},
		"when the store fails": {
			id:    "task-1",
			input: usecase.CompleteInput{Completed: true},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					SetCompleted(gomock.Any(), "task-1", true).
					Return(nil, apperrors.NewApplicationError(mappings.TaskCompleteStoreError, errors.New("permission denied")))
			},
			expectedErr:  mappings.TaskCompleteStoreError,
			expectsError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repository: mock_repository.NewMockRepository(ctrl)}
			tt.prepare(&f)

			output, appErr := usecase.NewCompleteUsecase(newContextFactory(f.repository)).
				Execute(context.Background(), tt.id, tt.input)

			if tt.expectsError {
				assert.Nil(t, output)
				assert.NotNil(t, appErr)
				assert.True(t, appErr.IsBasedOn(tt.expectedErr))

				return
			}

			assert.Nil(t, appErr)
			assert.Equal(t, tt.expectedCompleted, output.Data.Completed)
		})
	}
}
