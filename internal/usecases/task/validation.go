package task

import (
	"errors"

	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

const (
	taskTitleMinLength       = 3
	taskTitleMaxLength       = 120
	taskDescriptionMaxLength = 500
)

func validate(task domain.Task) apperrors.ApplicationError {
	switch {
	case task.Title == "":
		return apperrors.NewApplicationError(mappings.TaskTitleRequiredError, errors.New("title is required"))
	case len([]rune(task.Title)) < taskTitleMinLength:
		return apperrors.NewApplicationError(mappings.TaskTitleTooShortError, errors.New("title is too short"))
	case len([]rune(task.Title)) > taskTitleMaxLength:
		return apperrors.NewApplicationError(mappings.TaskTitleTooLongError, errors.New("title is too long"))
	case len([]rune(task.Description)) > taskDescriptionMaxLength:
		return apperrors.NewApplicationError(mappings.TaskDescriptionTooLongError, errors.New("description is too long"))
	}

	return nil
}
