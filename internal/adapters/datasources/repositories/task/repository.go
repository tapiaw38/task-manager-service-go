package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/jsonstore"
	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

type (
	Repository interface {
		Create(context.Context, domain.Task) (*domain.Task, apperrors.ApplicationError)
		List(context.Context) ([]domain.Task, int, apperrors.ApplicationError)
		Get(context.Context, string) (*domain.Task, apperrors.ApplicationError)
		SetCompleted(context.Context, string, bool) (*domain.Task, apperrors.ApplicationError)
		Delete(context.Context, string) apperrors.ApplicationError
	}

	repository struct {
		store *jsonstore.Store
		now   func() time.Time
		newID func() string
	}
)

func NewRepository(store *jsonstore.Store) Repository {
	return &repository{
		store: store,
		now:   func() time.Time { return time.Now().UTC() },
		newID: func() string { return uuid.NewString() },
	}
}

func findIndexByID(data []domain.Task, id string) int {
	for i := range data {
		if data[i].ID == id {
			return i
		}
	}

	return -1
}

func notFoundError(id string) apperrors.ApplicationError {
	return apperrors.NewApplicationError(mappings.TaskNotFoundError, nil).
		AddExtraFields(map[string]any{"task_id": id})
}

func toApplicationError(err error, fallback mappings.ErrorDetails) apperrors.ApplicationError {
	if appErr, ok := err.(apperrors.ApplicationError); ok {
		return appErr
	}

	return apperrors.NewApplicationError(fallback, err)
}
