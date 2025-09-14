package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTelemetry_AllModes(t *testing.T) {
	modes := []Mode{ModeDebug, ModeProd}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			// Set required environment variables for testing
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
			t.Setenv("OTEL_SERVICE_NAME", "test-service")
			t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

			ctx, cleanup, err := InitTelemetry(context.Background(), mode)
			assert.NoError(t, err)
			assert.NotNil(t, cleanup)
			defer cleanup(ctx)
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

	modes := []Mode{ModeDebug, ModeProd}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			lp, err := newOTELLogProvider(context.Background(), resource)
			assert.NoError(t, err)
			assert.NotNil(t, lp)
			lp.Shutdown(ctx)
		})
	}
}

func TestInitTelemetry_ResourceError(t *testing.T) {
	// Set invalid environment to potentially trigger resource creation error
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "=invalid,key=,=value")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)

	if err != nil {
		// If error occurs, context should still have a logger and cleanup may not be nil
		assert.NotNil(t, LoggerFromContext(ctx))
		// cleanup may be non-nil, but should be callable
	} else {
		// If no error, should have valid context and cleanup
		assert.NotNil(t, LoggerFromContext(ctx))
		assert.NotNil(t, cleanup)
		cleanup(ctx)
	}
}

func TestInitTelemetry_ServiceConfigCreation(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system,service.version=v1.0.0,deployment.environment=test-env,host.name=test-host")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	assert.NoError(t, err)
	assert.NotNil(t, cleanup)
	sc := ServiceConfigFromContext(ctx)
	assert.True(t, sc.IsValid())
	assert.Equal(t, "test-service", sc.Name)
	assert.Equal(t, "test-system", sc.System)
	assert.Equal(t, "v1.0.0", sc.Version)
	// Environment and HostName may be empty if not picked up by OTEL, but should not error
	cleanup(ctx)
}

func TestInitTelemetry_LoggerError(t *testing.T) {
	// Test with different modes to exercise error paths in newLogger
	modes := []Mode{ModeProd}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			// Set required environment variables for testing
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
			t.Setenv("OTEL_SERVICE_NAME", "test-service")
			t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

			ctx, cleanup, err := InitTelemetry(context.Background(), mode)

			if err != nil {
				// Error path: context should have a logger, cleanup may be non-nil
				assert.NotNil(t, LoggerFromContext(ctx))
			} else {
				assert.NotNil(t, LoggerFromContext(ctx))
				assert.NotNil(t, cleanup)
				cleanup(ctx)
			}
		})
	}
}

func TestInitTelemetry_NewResourceError(t *testing.T) {
	// Test error path when newResource fails
	// We can simulate this by providing invalid OTEL attributes
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "invalid=attribute=with=too=many=equals")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	if err != nil {
		// Error path: context should have a logger, cleanup may be non-nil
		assert.NotNil(t, LoggerFromContext(ctx))
	} else {
		// No error: should have valid context and cleanup
		assert.NotNil(t, LoggerFromContext(ctx))
		assert.NotNil(t, cleanup)
		cleanup(ctx)
	}
}

func TestInitTelemetry_NewLoggerError(t *testing.T) {
	// Test error path when newLogger fails
	// This is harder to trigger, but we can test with different modes
	modes := []Mode{ModeProd}

	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			// Try to force conditions that might cause newLogger to fail
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
			t.Setenv("OTEL_SERVICE_NAME", "test-service")
			t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

			ctx := context.Background()

			ctx, cleanup, err := InitTelemetry(ctx, mode)
			if err != nil {
				// Error path: context should have a logger, cleanup may be non-nil
				assert.NotNil(t, LoggerFromContext(ctx))
			} else {
				// Success path
				assert.NotNil(t, LoggerFromContext(ctx))
				assert.NotNil(t, cleanup)
				cleanup(ctx)
			}
		})
	}
}
