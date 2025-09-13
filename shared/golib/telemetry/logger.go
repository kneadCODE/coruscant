package telemetry

import (
	"context"
	"log/slog"
)

// contextKey is a private type used as a key for storing logger in context
type contextKey struct{}

// loggerKey is the key used to store logger in context
var loggerKey = contextKey{}

// SetLoggerFieldsInContext returns a new context with a logger that has additional fields.
// This is useful for adding request-scoped fields like request_id, user_id, etc.
func SetLoggerFieldsInContext(ctx context.Context, args ...any) context.Context {
	logger := LoggerFromContext(ctx)
	if logger == nil {
		return ctx
	}
	return setLoggerInContext(ctx, logger.With(args...))
}

// LoggerFromContext retrieves the logger from the context, or returns nil if not present.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return nil
}

// setLoggerInContext returns a new context with the given logger attached.
// This is typically called at the beginning of request handling to attach request-specific logging context.
func setLoggerInContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
