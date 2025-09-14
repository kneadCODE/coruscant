package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
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

// RecordErrorEvent logs an error-level event if a logger is present in the context.
// It also adds the error event to the active span if tracing is enabled.
func RecordErrorEvent(ctx context.Context, err error, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.ErrorContext(ctx, err.Error(), args...)
	}

	// Record error event in active span if present
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.RecordError(err, trace.WithTimestamp(time.Now()))
		addSpanEvent(ctx, "error occurred", args...)
	}
}

// addSpanEvent adds an event to the active span with the given message and attributes.
// This consolidates the repeated span event logic used across all record functions.
func addSpanEvent(ctx context.Context, message string, args ...any) {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		attrs := make([]attribute.KeyValue, 0, len(args)/2)
		for i := 0; i < len(args)-1; i += 2 {
			if key, ok := args[i].(string); ok {
				attrs = append(attrs, attribute.String(key, formatValue(args[i+1])))
			}
		}
		span.AddEvent(message, trace.WithTimestamp(time.Now()), trace.WithAttributes(attrs...))
	}
}

// formatValue converts any value to a string representation suitable for tracing attributes
func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%.2f", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
