package telemetry

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"
)

func TestRecordDebugEvent_WithLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))
	ctx := setLoggerInContext(context.Background(), logger)
	// Should not panic or error
	RecordDebugEvent(ctx, "debug", "foo", "bar")
}

func TestRecordDebugEvent_WithoutLogger(t *testing.T) {
	ctx := context.Background()
	// Should not panic or error
	RecordDebugEvent(ctx, "debug", "foo", "bar")
}

func TestRecordInfoEvent_WithLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))
	ctx := setLoggerInContext(context.Background(), logger)
	// Should not panic or error
	RecordInfoEvent(ctx, "info", "foo", "bar")
}

func TestRecordInfoEvent_WithoutLogger(t *testing.T) {
	ctx := context.Background()
	// Should not panic or error
	RecordInfoEvent(ctx, "info", "foo", "bar")
}

func TestRecordErrorEvent_WithLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))
	ctx := setLoggerInContext(context.Background(), logger)
	RecordErrorEvent(ctx, errors.New("fail"), "foo", "bar")
}

func TestRecordErrorEvent_WithoutLogger(t *testing.T) {
	ctx := context.Background()
	RecordErrorEvent(ctx, errors.New("fail"), "foo", "bar")
}
