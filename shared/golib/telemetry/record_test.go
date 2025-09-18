package telemetry

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordEvents(t *testing.T) {
	tests := []struct {
		name            string
		withLogger      bool
		logLevel        slog.Level
		recordFunc      func(context.Context)
		expectedMessage string
	}{
		{
			name:       "RecordDebugEvent with logger",
			withLogger: true,
			logLevel:   slog.LevelDebug,
			recordFunc: func(ctx context.Context) {
				RecordDebugEvent(ctx, "debug message", "foo", "bar")
			},
			expectedMessage: "debug message",
		},
		{
			name:       "RecordDebugEvent without logger",
			withLogger: false,
			recordFunc: func(ctx context.Context) {
				RecordDebugEvent(ctx, "debug message", "foo", "bar")
			},
		},
		{
			name:       "RecordInfoEvent with logger",
			withLogger: true,
			logLevel:   slog.LevelInfo,
			recordFunc: func(ctx context.Context) {
				RecordInfoEvent(ctx, "info message", "foo", "bar")
			},
			expectedMessage: "info message",
		},
		{
			name:       "RecordInfoEvent without logger",
			withLogger: false,
			recordFunc: func(ctx context.Context) {
				RecordInfoEvent(ctx, "info message", "foo", "bar")
			},
		},
		{
			name:       "RecordErrorEvent with logger",
			withLogger: true,
			logLevel:   slog.LevelError,
			recordFunc: func(ctx context.Context) {
				RecordErrorEvent(ctx, errors.New("test error message"), "foo", "bar")
			},
			expectedMessage: "test error message",
		},
		{
			name:       "RecordErrorEvent without logger",
			withLogger: false,
			recordFunc: func(ctx context.Context) {
				RecordErrorEvent(ctx, errors.New("test error message"), "foo", "bar")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			var buf *bytes.Buffer

			if tt.withLogger {
				buf = &bytes.Buffer{}
				handlerOpts := &slog.HandlerOptions{Level: tt.logLevel}
				logger := slog.New(slog.NewTextHandler(buf, handlerOpts))
				ctx = setLoggerInContext(context.Background(), logger)
			} else {
				ctx = context.Background()
			}

			// Test that the record function doesn't panic
			assert.NotPanics(t, func() {
				tt.recordFunc(ctx)
			})

			// If we have a logger, verify the message was logged
			if tt.withLogger && tt.expectedMessage != "" {
				assert.Contains(t, buf.String(), tt.expectedMessage)
			}
		})
	}
}
