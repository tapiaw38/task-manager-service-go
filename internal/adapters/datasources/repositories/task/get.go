package task

import (
	"context"

	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Task, apperrors.ApplicationError) {
	if err := ctx.Err(); err != nil {
		return nil, apperrors.NewApplicationError(mappings.TaskGetStoreError, err)
	}

	data := r.store.Snapshot()

	index := findIndexByID(data, id)
	if index < 0 {
		return nil, notFoundError(id)
	}

	found := data[index]

	return &found, nil
}
