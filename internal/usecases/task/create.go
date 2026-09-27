package task

import (
	"context"

	"github.com/tapiaw38/task-manager-service-go/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
)

type (
	CreateUsecase interface {
		Execute(context.Context, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateOutput struct {
		Data TaskOutputData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{
		contextFactory: contextFactory,
	}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	task := input.toDomain()

	if appErr := validate(task); appErr != nil {
		return nil, appErr
	}

	created, appErr := app.Repositories.Task.Create(ctx, task)
	if appErr != nil {
		return nil, appErr
	}

	return &CreateOutput{Data: toTaskOutputData(*created)}, nil
}
