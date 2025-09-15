package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// RecordDebugEvent logs a debug-level event if a logger is present in the context.
// It also adds the event to the active span if tracing is enabled.
func RecordDebugEvent(ctx context.Context, msg string, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.DebugContext(ctx, msg, args...)
	}

	// Record event in active span if present with timestamp
	addSpanEvent(ctx, msg, args...)
}

// RecordInfoEvent logs an info-level event if a logger is present in the context.
// It also adds the event to the active span if tracing is enabled.
func RecordInfoEvent(ctx context.Context, message string, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.InfoContext(ctx, message, args...)
	}

	// Record event in active span if present with timestamp
	addSpanEvent(ctx, message, args...)
}

// RecordWarnEvent logs an warn-level event if a logger is present in the context.
// It also adds the event to the active span if tracing is enabled.
func RecordWarnEvent(ctx context.Context, message string, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.WarnContext(ctx, message, args...)
	}

	// Record event in active span if present with timestamp
	addSpanEvent(ctx, message, args...)
}

// RecordErrorEvent logs an error-level event if a logger is present in the context.
// It also adds the error event to the active span if tracing is enabled.
func RecordErrorEvent(ctx context.Context, err error, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.ErrorContext(ctx, err.Error(), args...)
	}

	// Record error event in active span if present
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.RecordError(err, trace.WithTimestamp(time.Now()))
		addSpanEvent(ctx, fmt.Sprintf("error occurred: %s", err.Error()), args...)
	}
}

// addSpanEvent adds an event to the active span with the given message and attributes.
// This consolidates the repeated span event logic used across all record functions.
func addSpanEvent(ctx context.Context, message string, args ...any) {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.AddEvent(message, trace.WithTimestamp(time.Now()), trace.WithAttributes(convertToOTELAttributes(args)...))
	}
}
