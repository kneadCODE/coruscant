package telemetry

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordDebugEvent_WithLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	ctx := setLoggerInContext(context.Background(), logger)

	assert.NotPanics(t, func() {
		RecordDebugEvent(ctx, "debug message", "foo", "bar")
	})
	assert.Contains(t, buf.String(), "debug message")
}

func TestRecordDebugEvent_WithoutLogger(t *testing.T) {
	ctx := context.Background()

	assert.NotPanics(t, func() {
		RecordDebugEvent(ctx, "debug message", "foo", "bar")
	})
}

func TestRecordInfoEvent_WithLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))
	ctx := setLoggerInContext(context.Background(), logger)

	assert.NotPanics(t, func() {
		RecordInfoEvent(ctx, "info message", "foo", "bar")
	})
	assert.Contains(t, buf.String(), "info message")
}

func TestRecordInfoEvent_WithoutLogger(t *testing.T) {
	ctx := context.Background()

	assert.NotPanics(t, func() {
		RecordInfoEvent(ctx, "info message", "foo", "bar")
	})
}

func TestRecordErrorEvent_WithLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))
	ctx := setLoggerInContext(context.Background(), logger)

	assert.NotPanics(t, func() {
		RecordErrorEvent(ctx, errors.New("test error message"), "foo", "bar")
	})
	assert.Contains(t, buf.String(), "test error message")
}

func TestRecordErrorEvent_WithoutLogger(t *testing.T) {
	ctx := context.Background()

	assert.NotPanics(t, func() {
		RecordErrorEvent(ctx, errors.New("test error message"), "foo", "bar")
	})
}
