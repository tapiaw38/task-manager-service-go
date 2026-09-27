package task

import (
	"context"

	"github.com/tapiaw38/task-manager-service-go/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
)

type (
	ListUsecase interface {
		Execute(context.Context) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data  []TaskOutputData `json:"data"`
		Total int              `json:"total"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{
		contextFactory: contextFactory,
	}
}

func (u *listUsecase) Execute(ctx context.Context) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	tasks, total, appErr := app.Repositories.Task.List(ctx)
	if appErr != nil {
		return nil, appErr
	}

	data := make([]TaskOutputData, 0, len(tasks))
	for _, task := range tasks {
		data = append(data, toTaskOutputData(task))
	}

	return &ListOutput{Data: data, Total: total}, nil
}
