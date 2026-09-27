package task

import (
	"context"

	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

func (r *repository) Create(ctx context.Context, t domain.Task) (*domain.Task, apperrors.ApplicationError) {
	if err := ctx.Err(); err != nil {
		return nil, apperrors.NewApplicationError(mappings.TaskCreateStoreError, err)
	}

	created := t
	created.ID = r.newID()
	created.CreatedAt = r.now()
	created.UpdatedAt = created.CreatedAt

	err := r.store.Mutate(func(data []domain.Task) ([]domain.Task, error) {
		return append(data, created), nil
	})
	if err != nil {
		return nil, toApplicationError(err, mappings.TaskCreateStoreError)
	}

	return &created, nil
}
