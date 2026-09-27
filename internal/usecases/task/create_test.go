package task_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	mock_repository "github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/repositories/task/mocks"
	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/task-manager-service-go/internal/usecases/task"
	"go.uber.org/mock/gomock"
)

func TestCreateUsecase(t *testing.T) {
	type fields struct {
		repository *mock_repository.MockRepository
	}

	tests := map[string]struct {
		input        usecase.CreateInput
		prepare      func(f *fields)
		expectedTask usecase.TaskOutputData
		expectedErr  mappings.ErrorDetails
		expectsError bool
	}{
		"when the task is created it sanitizes the input": {
			input: usecase.CreateInput{
				Title:       "  Comprar   leche ",
				Description: "  Ir al  supermercado ",
			},
			prepare: func(f *fields) {
				created := newTask("task-1", "Comprar leche")
				f.repository.EXPECT().
					Create(gomock.Any(), domain.Task{
						Title:       "Comprar leche",
						Description: "Ir al supermercado",
					}).
					Return(&created, nil)
			},
			expectedTask: usecase.TaskOutputData{
				ID:          "task-1",
				Title:       "Comprar leche",
				Description: "Description of Comprar leche",
				Completed:   false,
				CreatedAt:   referenceTime,
				UpdatedAt:   referenceTime,
			},
		},
		"when the title is empty": {
			input:        usecase.CreateInput{Title: "   ", Description: "Something"},
			prepare:      func(f *fields) {},
			expectedErr:  mappings.TaskTitleRequiredError,
			expectsError: true,
		},
		"when the title is too short": {
			input:        usecase.CreateInput{Title: "ab", Description: "Something"},
			prepare:      func(f *fields) {},
			expectedErr:  mappings.TaskTitleTooShortError,
			expectsError: true,
		},
		"when the title is too long": {
			input:        usecase.CreateInput{Title: strings.Repeat("a", 121), Description: "Something"},
			prepare:      func(f *fields) {},
			expectedErr:  mappings.TaskTitleTooLongError,
			expectsError: true,
		},
		"when the description is too long": {
			input:        usecase.CreateInput{Title: "Valid title", Description: strings.Repeat("a", 501)},
			prepare:      func(f *fields) {},
			expectedErr:  mappings.TaskDescriptionTooLongError,
			expectsError: true,
		},
		"when the store fails": {
			input: usecase.CreateInput{Title: "Valid title", Description: "Something"},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.NewApplicationError(mappings.TaskCreateStoreError, errors.New("disk is full")))
			},
			expectedErr:  mappings.TaskCreateStoreError,
			expectsError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repository: mock_repository.NewMockRepository(ctrl)}
			tt.prepare(&f)

			output, appErr := usecase.NewCreateUsecase(newContextFactory(f.repository)).
				Execute(context.Background(), tt.input)

			if tt.expectsError {
				assert.Nil(t, output)
				assert.NotNil(t, appErr)
				assert.True(t, appErr.IsBasedOn(tt.expectedErr))

				return
			}

			assert.Nil(t, appErr)
			assert.Equal(t, tt.expectedTask, output.Data)
		})
	}
}
