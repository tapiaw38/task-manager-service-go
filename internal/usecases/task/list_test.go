package task_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	mock_repository "github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/repositories/task/mocks"
	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/task-manager-service-go/internal/usecases/task"
	"go.uber.org/mock/gomock"
)

func TestListUsecase(t *testing.T) {
	type fields struct {
		repository *mock_repository.MockRepository
	}

	tests := map[string]struct {
		prepare       func(f *fields)
		expectedIDs   []string
		expectedTotal int
		expectedErr   mappings.ErrorDetails
		expectsError  bool
	}{
		"when there are tasks it maps them to the output contract": {
			prepare: func(f *fields) {
				f.repository.EXPECT().
					List(gomock.Any()).
					Return([]domain.Task{
						newTask("task-1", "First task"),
						newTask("task-2", "Second task"),
					}, 2, nil)
			},
			expectedIDs:   []string{"task-1", "task-2"},
			expectedTotal: 2,
		},
		"when there are no tasks it returns an empty list": {
			prepare: func(f *fields) {
				f.repository.EXPECT().
					List(gomock.Any()).
					Return([]domain.Task{}, 0, nil)
			},
			expectedIDs:   []string{},
			expectedTotal: 0,
		},
		"when the store fails": {
			prepare: func(f *fields) {
				f.repository.EXPECT().
					List(gomock.Any()).
					Return(nil, 0, apperrors.NewApplicationError(mappings.TaskListStoreError, errors.New("corrupted file")))
			},
			expectedErr:  mappings.TaskListStoreError,
			expectsError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repository: mock_repository.NewMockRepository(ctrl)}
			tt.prepare(&f)

			output, appErr := usecase.NewListUsecase(newContextFactory(f.repository)).
				Execute(context.Background())

			if tt.expectsError {
				assert.Nil(t, output)
				assert.NotNil(t, appErr)
				assert.True(t, appErr.IsBasedOn(tt.expectedErr))

				return
			}

			assert.Nil(t, appErr)
			assert.Equal(t, tt.expectedTotal, output.Total)

			ids := make([]string, 0, len(output.Data))
			for _, item := range output.Data {
				ids = append(ids, item.ID)
			}

			assert.Equal(t, tt.expectedIDs, ids)
		})
	}
}
