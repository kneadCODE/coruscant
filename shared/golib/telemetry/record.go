package telemetry

import (
	"context"
)

// RecordDebugEvent logs a debug-level event if a logger is present in the context.
func RecordDebugEvent(ctx context.Context, msg string, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.DebugContext(ctx, msg, args...)
	}

	// TODO: Check if span exists and then record event in there
}

// RecordInfoEvent logs an info-level event if a logger is present in the context.
func RecordInfoEvent(ctx context.Context, message string, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.InfoContext(ctx, message, args...)
	}

	// TODO: Check if span exists and then record event in there
}

// RecordErrorEvent logs an error-level event if a logger is present in the context.
func RecordErrorEvent(ctx context.Context, err error, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.ErrorContext(ctx, err.Error(), args...)
	}

	// TODO: Check if span exists and then record event in there
}
