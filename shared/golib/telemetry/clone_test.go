package telemetry

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneCopiesLogger(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	ctx := setLoggerInContext(context.Background(), logger)
	cloned := Clone(ctx)
	assert.Equal(t, LoggerFromContext(ctx), LoggerFromContext(cloned))
}

func TestCloneEmptyContext(t *testing.T) {
	ctx := context.Background()
	cloned := Clone(ctx)
	assert.Nil(t, LoggerFromContext(cloned))
	assert.Nil(t, MetricsCollectorFromContext(cloned))
}

func TestCloneCopiesMetricsCollector(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_SERVICE_NAME", "test-service")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	require.NoError(t, err)
	defer cleanup(ctx)

	collector, err := NewMetricsCollector()
	require.NoError(t, err)

	ctx = setMetricsCollectorInContext(ctx, collector)
	cloned := Clone(ctx)

	assert.Equal(t, MetricsCollectorFromContext(ctx), MetricsCollectorFromContext(cloned))
	assert.NotNil(t, MetricsCollectorFromContext(cloned))
}

func TestCloneCopiesBothLoggerAndMetricsCollector(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_SERVICE_NAME", "test-service")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	require.NoError(t, err)
	defer cleanup(ctx)

	logger := slog.New(slog.NewTextHandler(nil, nil))
	collector, err := NewMetricsCollector()
	require.NoError(t, err)

	ctx = setLoggerInContext(ctx, logger)
	ctx = setMetricsCollectorInContext(ctx, collector)
	cloned := Clone(ctx)

	assert.Equal(t, LoggerFromContext(ctx), LoggerFromContext(cloned))
	assert.Equal(t, MetricsCollectorFromContext(ctx), MetricsCollectorFromContext(cloned))
	assert.NotNil(t, LoggerFromContext(cloned))
	assert.NotNil(t, MetricsCollectorFromContext(cloned))
}
