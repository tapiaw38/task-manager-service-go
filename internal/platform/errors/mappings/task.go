package mappings

import "net/http"

var (
	TaskTitleRequiredError = ErrorDetails{
		"task:validation:title-required",
		http.StatusBadRequest,
		"task title is required",
	}

	TaskTitleTooShortError = ErrorDetails{
		"task:validation:title-too-short",
		http.StatusBadRequest,
		"task title is too short",
	}

	TaskTitleTooLongError = ErrorDetails{
		"task:validation:title-too-long",
		http.StatusBadRequest,
		"task title is too long",
	}

	TaskDescriptionTooLongError = ErrorDetails{
		"task:validation:description-too-long",
		http.StatusBadRequest,
		"task description is too long",
	}

	TaskNotFoundError = ErrorDetails{
		"task:shared:not-found",
		http.StatusNotFound,
		"task not found",
	}

	TaskCreateStoreError = ErrorDetails{
		"task:create:store-error",
		http.StatusInternalServerError,
		"failed to create task",
	}

	TaskListStoreError = ErrorDetails{
		"task:list:store-error",
		http.StatusInternalServerError,
		"failed to list tasks",
	}

	TaskGetStoreError = ErrorDetails{
		"task:get:store-error",
		http.StatusInternalServerError,
		"failed to retrieve task",
	}

	TaskCompleteStoreError = ErrorDetails{
		"task:complete:store-error",
		http.StatusInternalServerError,
		"failed to update task status",
	}

	TaskDeleteStoreError = ErrorDetails{
		"task:delete:store-error",
		http.StatusInternalServerError,
		"failed to delete task",
	}
)
