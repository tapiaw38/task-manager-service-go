package task

import (
	"context"

	"github.com/tapiaw38/task-manager-service-go/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
)

type (
	CompleteUsecase interface {
		Execute(context.Context, string, CompleteInput) (*CompleteOutput, apperrors.ApplicationError)
	}

	completeUsecase struct {
		contextFactory appcontext.Factory
	}

	CompleteOutput struct {
		Data TaskOutputData `json:"data"`
	}
)

func NewCompleteUsecase(contextFactory appcontext.Factory) CompleteUsecase {
	return &completeUsecase{
		contextFactory: contextFactory,
	}
}

func (u *completeUsecase) Execute(ctx context.Context, id string, input CompleteInput) (*CompleteOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	updated, appErr := app.Repositories.Task.SetCompleted(ctx, id, input.Completed)
	if appErr != nil {
		return nil, appErr
	}

	return &CompleteOutput{Data: toTaskOutputData(*updated)}, nil
}
