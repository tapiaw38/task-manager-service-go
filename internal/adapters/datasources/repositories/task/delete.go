package task

import (
	"context"

	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

func (r *repository) Delete(ctx context.Context, id string) apperrors.ApplicationError {
	if err := ctx.Err(); err != nil {
		return apperrors.NewApplicationError(mappings.TaskDeleteStoreError, err)
	}

	err := r.store.Mutate(func(data []domain.Task) ([]domain.Task, error) {
		index := findIndexByID(data, id)
		if index < 0 {
			return nil, notFoundError(id)
		}

		return append(data[:index], data[index+1:]...), nil
	})
	if err != nil {
		return toApplicationError(err, mappings.TaskDeleteStoreError)
	}

	return nil
}
