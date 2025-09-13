package telemetry

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeasureBasic(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test basic span creation with Measure
	ctx, finishFunc := Measure(ctx, "test-operation")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestMeasureWithStringPairs(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test measure with string key-value pairs
	ctx, finishFunc := Measure(ctx, "test-operation", "test.key", "test.value", "test.number", "42")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestMeasureServerOperation(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	ctx, finishFunc := Measure(ctx, "http-request", "operation.type", "server")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestMeasureClientOperation(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	ctx, finishFunc := Measure(ctx, "database-query", "operation.type", "client")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestRecordSuccessEvent(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	ctx, finishFunc := Measure(ctx, "test-operation")
	defer finishFunc(nil)

	// Test recording info event (success message)
	RecordInfoEvent(ctx, "Operation completed successfully", "result", "success", "duration", 100)

	// Should not panic or cause issues
	assert.NotNil(t, ctx)
}

func TestRecordErrorEventInMeasure(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	ctx, finishFunc := Measure(ctx, "test-operation")
	defer finishFunc(nil)

	testError := errors.New("test error")
	RecordErrorEvent(ctx, testError, "step", "validation")

	assert.NotNil(t, ctx)
}

func TestRecordInfoEventInMeasure(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	ctx, finishFunc := Measure(ctx, "test-operation")
	defer finishFunc(nil)

	RecordInfoEvent(ctx, "Processing started", "step", "initialization", "count", 10)

	assert.NotNil(t, ctx)
}

func TestRecordDebugEventInMeasure(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	ctx, finishFunc := Measure(ctx, "test-operation")
	defer finishFunc(nil)

	RecordDebugEvent(ctx, "Debug information", "variable", "value", "state", "active")

	assert.NotNil(t, ctx)
}

func TestMeasureWithMultipleAttributes(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test measure with multiple attributes
	ctx, finishFunc := Measure(ctx, "test-operation", "operation.type", "test", "operation.success", "true", "user.id", "123")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestMeasureSpanContext(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test measure creating span context
	ctx, finishFunc := Measure(ctx, "test-operation")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Test that we can create nested measures
	ctx2, finishFunc2 := Measure(ctx, "nested-operation")
	assert.NotNil(t, ctx2)
	assert.NotNil(t, finishFunc2)

	// Finish in reverse order
	finishFunc2(nil)
	finishFunc(nil)
}

func TestMeasureTelemetryInitialization(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test that telemetry is initialized and measure works
	ctx, finishFunc := Measure(ctx, "test-operation")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestMeasureNestingWithEvents(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Parent measure
	ctx, parentFinish := Measure(ctx, "parent-operation")
	RecordInfoEvent(ctx, "Parent operation started")

	// Child measure
	ctx, childFinish := Measure(ctx, "child-operation")
	RecordInfoEvent(ctx, "Child operation started")

	// Finish child then parent
	childFinish(nil)
	parentFinish(nil)

	// Should not panic
	assert.NotNil(t, ctx)
}

func TestMeasureErrorHandling(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	ctx, finishFunc := Measure(ctx, "test-operation")

	// Simulate an error occurring
	testError := errors.New("something went wrong")
	RecordErrorEvent(ctx, testError, "component", "validator")

	// Finish with error to record it in the span
	finishFunc(testError)

	// Should not panic
	assert.NotNil(t, ctx)
}

// Benchmark tests for performance validation
func BenchmarkMeasureOperation(b *testing.B) {
	// Set required environment variables for testing
	b.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	b.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(b, err)
	defer cleanup(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, finishFunc := Measure(ctx, "benchmark-operation")
		finishFunc(nil)
	}
}

func TestMeasure(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test basic measure without attributes
	ctx, finishFunc := Measure(ctx, "test-operation")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestMeasureWithAttributes(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test measure with different attribute types
	ctx, finishFunc := Measure(ctx, "test-operation-with-attrs",
		"user.id", "12345",
		"operation.type", "read",
		"retry.count", 3,
		"duration.ms", 150.75,
		"success", true,
	)
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestMeasureWithError(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test measure with error
	ctx, finishFunc := Measure(ctx, "test-operation-with-error", "component", "database")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish with error
	testError := errors.New("database connection failed")
	finishFunc(testError)
}

func TestMeasureWithOddNumberOfAttributes(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test measure with odd number of attributes (should be ignored)
	ctx, finishFunc := Measure(ctx, "test-operation-odd-attrs", "user.id", "12345", "incomplete-key")
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestMeasureWithInvalidKeyTypes(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test measure with non-string keys (should be skipped)
	ctx, finishFunc := Measure(ctx, "test-operation-invalid-keys",
		123, "invalid-key", // invalid key type (int)
		"valid.key", "valid-value",
		true, "another-invalid", // invalid key type (bool)
		"custom.type", struct{ Name string }{Name: "test"}, // custom type as value
	)
	assert.NotNil(t, ctx)
	assert.NotNil(t, finishFunc)

	// Finish without error
	finishFunc(nil)
}

func TestMeasureNested(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(t, err)
	defer cleanup(ctx)

	// Test nested measure operations
	ctx, parentFinish := Measure(ctx, "parent-operation", "level", "1")
	assert.NotNil(t, ctx)

	ctx, childFinish := Measure(ctx, "child-operation", "level", "2", "parent", "parent-operation")
	assert.NotNil(t, ctx)

	// Finish child first
	childFinish(nil)

	// Then finish parent
	parentFinish(nil)
}

func BenchmarkMeasure(b *testing.B) {
	// Set required environment variables for testing
	b.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	b.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(b, err)
	defer cleanup(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, finishFunc := Measure(ctx, "benchmark-operation", "iteration", "value")
		finishFunc(nil)
	}
}

func BenchmarkRecordInfoEvent(b *testing.B) {
	// Set required environment variables for testing
	b.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	b.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	ctx, cleanup, err := InitTelemetry(ctx, ModeDev)
	require.NoError(b, err)
	defer cleanup(ctx)

	ctx, finishFunc := Measure(ctx, "benchmark-span")
	defer finishFunc(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RecordInfoEvent(ctx, "benchmark message", "iteration", i)
	}
}
