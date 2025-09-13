package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// RecordDebugEvent logs a debug-level event if a logger is present in the context.
// It also adds the event to the active span if tracing is enabled.
func RecordDebugEvent(ctx context.Context, msg string, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.DebugContext(ctx, msg, args...)
	}

	// Record event in active span if present
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.AddEvent(msg)
	}
}

// RecordInfoEvent logs an info-level event if a logger is present in the context.
// It also adds the event to the active span if tracing is enabled.
func RecordInfoEvent(ctx context.Context, message string, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.InfoContext(ctx, message, args...)
	}

	// Record event in active span if present
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.AddEvent(message)
	}
}

// RecordErrorEvent logs an error-level event if a logger is present in the context.
// It also adds the error event to the active span if tracing is enabled.
func RecordErrorEvent(ctx context.Context, err error, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.ErrorContext(ctx, err.Error(), args...)
	}

	// Record error event in active span if present
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.RecordError(err)
		span.AddEvent("error occurred")
	}
}
