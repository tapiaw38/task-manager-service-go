package task

import (
	"context"

	"github.com/tapiaw38/task-manager-service-go/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
)

type (
	GetUsecase interface {
		Execute(context.Context, string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data TaskOutputData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{
		contextFactory: contextFactory,
	}
}

func (u *getUsecase) Execute(ctx context.Context, id string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	task, appErr := app.Repositories.Task.Get(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	return &GetOutput{Data: toTaskOutputData(*task)}, nil
}
