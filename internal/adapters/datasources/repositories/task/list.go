package task

import (
	"context"
	"sort"

	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

func (r *repository) List(ctx context.Context) ([]domain.Task, int, apperrors.ApplicationError) {
	if err := ctx.Err(); err != nil {
		return nil, 0, apperrors.NewApplicationError(mappings.TaskListStoreError, err)
	}

	data := r.store.Snapshot()

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].CreatedAt.After(data[j].CreatedAt)
	})

	return data, len(data), nil
}
