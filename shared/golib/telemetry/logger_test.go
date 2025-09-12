package telemetry

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerContextAPI(t *testing.T) {
	logger := slog.Default()
	type fields struct {
		args []any
	}
	tests := []struct {
		name      string
		setup     func() context.Context
		want      *slog.Logger
		wantEqual bool
		wantNil   bool
		fields    fields
	}{
		{
			name: "SetLoggerInContext attaches logger",
			setup: func() context.Context {
				ctx := context.Background()
				return setLoggerInContext(ctx, logger)
			},
			want:      logger,
			wantEqual: true,
		},
		{
			name:      "LoggerFromContext returns nil if not set",
			setup:     context.Background,
			want:      nil,
			wantEqual: true,
		},
		{
			name: "SetLoggerFieldsInContext adds fields",
			setup: func() context.Context {
				ctx := context.Background()
				ctx = setLoggerInContext(ctx, logger)
				ctx = SetLoggerFieldsInContext(ctx, "request_id", "12345", "user_id", "user123")
				return ctx
			},
			wantNil: true,
			fields:  fields{args: []any{"request_id", "12345", "user_id", "user123"}},
		},
		{
			name:      "LoggerFromContext returns nil if not set",
			setup:     context.Background,
			wantNil:   false,
			want:      nil,
			wantEqual: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			got := LoggerFromContext(ctx)
			if tt.wantEqual {
				assert.Equal(t, tt.want, got)
			}
			if tt.wantNil {
				assert.NotNil(t, got)
			}
			// If fields are set, check logger has With applied (not default logger)
			if len(tt.fields.args) > 0 {
				// The logger returned should not be the same as slog.Default or the original logger
				assert.NotEqual(t, slog.Default(), got)
				assert.NotEqual(t, logger, got)
			}
		})
	}
}

func TestLoggerFromContextWithoutLogger(t *testing.T) {
	ctx := context.Background()
	logger := LoggerFromContext(ctx)
	assert.Nil(t, logger, "Expected nil when no logger in context")
}

func TestSetLoggerFieldsInContext(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()
	ctx = setLoggerInContext(ctx, logger)
	ctx = SetLoggerFieldsInContext(ctx, "request_id", "12345", "user_id", "user123")
	retrievedLogger := LoggerFromContext(ctx)
	assert.NotNil(t, retrievedLogger, "Expected non-nil logger")
}

func TestSetLoggerFieldsInContext_NilLogger(t *testing.T) {
	ctx := context.Background()
	resultCtx := SetLoggerFieldsInContext(ctx, "request_id", "12345")
	assert.Equal(t, ctx, resultCtx, "Should return same context when no logger present")
}
