package mappings

import "net/http"

var (
	InvalidRequestBodyError = ErrorDetails{
		"common:request-body-parsing-error",
		http.StatusBadRequest,
		"invalid request body format",
	}

	InvalidParamsError = ErrorDetails{
		"common:invalid-params",
		http.StatusBadRequest,
		"invalid request parameters",
	}

	InternalServerError = ErrorDetails{
		"common:internal-server-error",
		http.StatusInternalServerError,
		"an internal server error occurred",
	}

	NotFoundError = ErrorDetails{
		"common:not-found",
		http.StatusNotFound,
		"resource not found",
	}
)
