package task

import (
	"context"

	"github.com/tapiaw38/task-manager-service-go/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
)

type (
	DeleteUsecase interface {
		Execute(context.Context, string) apperrors.ApplicationError
	}

	deleteUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{
		contextFactory: contextFactory,
	}
}

func (u *deleteUsecase) Execute(ctx context.Context, id string) apperrors.ApplicationError {
	app := u.contextFactory()

	return app.Repositories.Task.Delete(ctx, id)
}
