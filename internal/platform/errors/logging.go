package errors

import (
	"context"
	"log/slog"
)

func (e *applicationError) Log(ctx context.Context) {
	attrs := []any{
		slog.String("internal_code", e.InternalCode()),
		slog.Int("status_code", e.StatusCode()),
		slog.String("message", e.Message()),
		slog.String("original_message", e.originalMessage),
	}

	for k, v := range e.extraFields {
		attrs = append(attrs, slog.Any(k, v))
	}

	slog.ErrorContext(ctx, "application error", attrs...)
}
