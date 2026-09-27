package errors

import (
	"context"
	"encoding/json"

	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

type ApplicationError interface {
	InternalCode() string
	StatusCode() int
	Message() string
	OriginalMessage() string
	AddExtraFields(fields map[string]any) ApplicationError
	IsBasedOn(mappings.ErrorDetails) bool
	Log(context.Context)
	json.Marshaler
	error
}

type applicationError struct {
	errorDetails    mappings.ErrorDetails
	originalMessage string
	extraFields     map[string]any
}

func (e *applicationError) InternalCode() string {
	return e.errorDetails.InternalCode
}

func (e *applicationError) StatusCode() int {
	return e.errorDetails.StatusCode
}

func (e *applicationError) Message() string {
	return e.errorDetails.Message
}

func (e *applicationError) OriginalMessage() string {
	return e.originalMessage
}

func (e *applicationError) Error() string {
	return e.errorDetails.Message
}

func (e *applicationError) AddExtraFields(fields map[string]any) ApplicationError {
	if e.extraFields == nil {
		e.extraFields = make(map[string]any)
	}

	for k, v := range fields {
		e.extraFields[k] = v
	}

	return e
}

func (e *applicationError) IsBasedOn(errorDetails mappings.ErrorDetails) bool {
	return e.errorDetails == errorDetails
}

func (e *applicationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(errorBody{
		Code:    e.InternalCode(),
		Message: e.Message(),
	})
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewApplicationError(details mappings.ErrorDetails, originalError error) ApplicationError {
	var originalMessage string

	if originalError != nil {
		originalMessage = originalError.Error()
	}

	return &applicationError{
		errorDetails:    details,
		originalMessage: originalMessage,
		extraFields:     make(map[string]any),
	}
}
