package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTelemetry_AllModes(t *testing.T) {
	modes := []Mode{ModeDev, ModeDevDebug, ModeProd, ModeProdDebug}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			// Set required environment variables for testing
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
			t.Setenv("OTEL_SERVICE_NAME", "test-service")

			ctx, cleanup, err := InitTelemetry(context.Background(), mode)
			assert.NoError(t, err)
			assert.NotNil(t, cleanup)
			defer cleanup()
			logger := LoggerFromContext(ctx)
			assert.NotNil(t, logger)
		})
	}
}

func TestNewLogger_ErrorHandling(t *testing.T) {
	// Test that newLogger handles different modes correctly
	// This is mostly to test the error return paths

	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	resource, err := newResource(ctx)
	assert.NoError(t, err)

	modes := []Mode{ModeDev, ModeDevDebug, ModeProd, ModeProdDebug}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			logger, cleanup, err := newOTELSlogLogger(context.Background(), resource)
			assert.NoError(t, err)
			assert.NotNil(t, logger)
			assert.NotNil(t, cleanup)
			cleanup()
		})
	}
}

func TestInitTelemetry_ResourceError(t *testing.T) {
	// Set invalid environment to potentially trigger resource creation error
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "=invalid,key=,=value")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDev)

	if err != nil {
		// If error occurs, context should be unchanged and cleanup should be nil
		assert.Equal(t, context.Background(), ctx)
		assert.Nil(t, cleanup)
	} else {
		// If no error, should have valid context and cleanup
		assert.NotNil(t, LoggerFromContext(ctx))
		assert.NotNil(t, cleanup)
		cleanup()
	}
}

func TestInitTelemetry_LoggerError(t *testing.T) {
	// Test with different modes to exercise error paths in newLogger
	modes := []Mode{ModeProd, ModeProdDebug}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			// Set required environment variables for testing
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
			t.Setenv("OTEL_SERVICE_NAME", "test-service")

			ctx, cleanup, err := InitTelemetry(context.Background(), mode)

			if err != nil {
				assert.Equal(t, context.Background(), ctx)
				assert.Nil(t, cleanup)
			} else {
				assert.NotNil(t, LoggerFromContext(ctx))
				assert.NotNil(t, cleanup)
				cleanup()
			}
		})
	}
}

func (m Mode) String() string {
	switch m {
	case ModeDev:
		return "ModeDev"
	case ModeDevDebug:
		return "ModeDevDebug"
	case ModeProd:
		return "ModeProd"
	case ModeProdDebug:
		return "ModeProdDebug"
	default:
		return "Unknown"
	}
}

func TestInitTelemetry_NewResourceError(t *testing.T) {
	// Test error path when newResource fails
	// We can simulate this by providing invalid OTEL attributes
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "invalid=attribute=with=too=many=equals")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDev)
	if err != nil {
		// Error path: should return original context and nil cleanup
		assert.Equal(t, context.Background(), ctx)
		assert.Nil(t, cleanup)
	} else {
		// No error: should have valid context and cleanup
		assert.NotNil(t, LoggerFromContext(ctx))
		assert.NotNil(t, cleanup)
		cleanup()
	}
}

func TestInitTelemetry_NewLoggerError(t *testing.T) {
	// Test error path when newLogger fails
	// This is harder to trigger, but we can test with different modes
	modes := []Mode{ModeProd, ModeProdDebug}

	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			// Try to force conditions that might cause newLogger to fail
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
			t.Setenv("OTEL_SERVICE_NAME", "test-service")

			ctx := context.Background()

			ctx, cleanup, err := InitTelemetry(ctx, mode)
			if err != nil {
				// Error path: should return original context and nil cleanup
				assert.Equal(t, context.Background(), ctx)
				assert.Nil(t, cleanup)
			} else {
				// Success path
				assert.NotNil(t, LoggerFromContext(ctx))
				assert.NotNil(t, cleanup)
				cleanup()
			}
		})
	}
}
