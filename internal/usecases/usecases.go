package usecases

import (
	"github.com/tapiaw38/task-manager-service-go/internal/platform/appcontext"
	"github.com/tapiaw38/task-manager-service-go/internal/usecases/task"
)

type Usecases struct {
	Task Task
}

type Task struct {
	CreateUsecase   task.CreateUsecase
	ListUsecase     task.ListUsecase
	GetUsecase      task.GetUsecase
	CompleteUsecase task.CompleteUsecase
	DeleteUsecase   task.DeleteUsecase
}

func CreateUsecases(contextFactory appcontext.Factory) *Usecases {
	return &Usecases{
		Task: Task{
			CreateUsecase:   task.NewCreateUsecase(contextFactory),
			ListUsecase:     task.NewListUsecase(contextFactory),
			GetUsecase:      task.NewGetUsecase(contextFactory),
			CompleteUsecase: task.NewCompleteUsecase(contextFactory),
			DeleteUsecase:   task.NewDeleteUsecase(contextFactory),
		},
	}
}
