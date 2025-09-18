package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/semconv/v1.37.0/dbconv"
)

func TestNewMetrics(t *testing.T) {
	ctx := context.Background()

	metrics, err := newMetrics(ctx, dbconv.SystemNamePostgreSQL, "localhost", 5432, "testdb")
	require.NoError(t, err)
	require.NotNil(t, metrics)

	assert.Equal(t, dbconv.SystemNamePostgreSQL, metrics.dbSystem)
	assert.NotEmpty(t, metrics.commonAttrs)
	assert.NotNil(t, metrics.operationDuration)
	assert.NotNil(t, metrics.returnedRows)
	assert.NotNil(t, metrics.connCount)
	assert.NotNil(t, metrics.connMax)
	assert.NotNil(t, metrics.connIdleMaxCount)
	assert.NotNil(t, metrics.connIdleMinCount)
	assert.NotNil(t, metrics.connPendingRequests)
	assert.NotNil(t, metrics.connTimetouts)
	assert.NotNil(t, metrics.connCreateTime)
	assert.NotNil(t, metrics.connWaitTime)
	assert.NotNil(t, metrics.connUseTime)
}

func TestMetrics_RecordConnectionStats(t *testing.T) {
	ctx := context.Background()
	metrics := createTestMetrics(t)

	// Test recording idle connections
	metrics.RecordConnectionStats(ctx, "idle", 5)

	// Test recording used connections
	metrics.RecordConnectionStats(ctx, "used", 3)

	// Should not panic or error
}

func TestMetrics_SetConnectionMax(t *testing.T) {
	ctx := context.Background()
	metrics := createTestMetrics(t)

	metrics.SetConnectionMax(ctx, 10)

	// Should not panic or error
}

func TestMetrics_SetConnectionIdleMax(t *testing.T) {
	ctx := context.Background()
	metrics := createTestMetrics(t)

	metrics.SetConnectionIdleMax(ctx, 8)

	// Should not panic or error
}

func TestMetrics_SetConnectionIdleMin(t *testing.T) {
	ctx := context.Background()
	metrics := createTestMetrics(t)

	metrics.SetConnectionIdleMin(ctx, 2)

	// Should not panic or error
}

func TestMetrics_RecordConnectionPendingRequests(t *testing.T) {
	ctx := context.Background()
	metrics := createTestMetrics(t)

	metrics.RecordConnectionPendingRequests(ctx, 1)

	// Should not panic or error
}

func TestMetrics_RecordConnectionTimeouts(t *testing.T) {
	ctx := context.Background()
	metrics := createTestMetrics(t)

	metrics.RecordConnectionTimeouts(ctx, 1)

	// Should not panic or error
}

func TestMetrics_RecordConnectionCreateTime(t *testing.T) {
	ctx := context.Background()
	metrics := createTestMetrics(t)

	duration := 100 * time.Millisecond
	metrics.RecordConnectionCreateTime(ctx, duration)

	// Should not panic or error
}

func TestMetrics_RecordConnectionWaitTime(t *testing.T) {
	ctx := context.Background()
	metrics := createTestMetrics(t)

	duration := 50 * time.Millisecond
	metrics.RecordConnectionWaitTime(ctx, duration)

	// Should not panic or error
}

func TestMetrics_RecordConnectionUseTime(t *testing.T) {
	ctx := context.Background()
	metrics := createTestMetrics(t)

	duration := 200 * time.Millisecond
	metrics.RecordConnectionUseTime(ctx, duration)

	// Should not panic or error
}

// Helper function to create test metrics
func createTestMetrics(t *testing.T) *Metrics {
	ctx := context.Background()
	metrics, err := newMetrics(ctx, dbconv.SystemNamePostgreSQL, "localhost", 5432, "testdb")
	require.NoError(t, err)
	return metrics
}
