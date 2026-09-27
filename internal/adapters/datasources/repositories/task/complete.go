package task

import (
	"context"

	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

func (r *repository) SetCompleted(ctx context.Context, id string, completed bool) (*domain.Task, apperrors.ApplicationError) {
	if err := ctx.Err(); err != nil {
		return nil, apperrors.NewApplicationError(mappings.TaskCompleteStoreError, err)
	}

	var updated domain.Task

	err := r.store.Mutate(func(data []domain.Task) ([]domain.Task, error) {
		index := findIndexByID(data, id)
		if index < 0 {
			return nil, notFoundError(id)
		}

		candidate := data[index]
		candidate.Completed = completed
		candidate.UpdatedAt = r.now()

		data[index] = candidate
		updated = candidate

		return data, nil
	})
	if err != nil {
		return nil, toApplicationError(err, mappings.TaskCompleteStoreError)
	}

	return &updated, nil
}
