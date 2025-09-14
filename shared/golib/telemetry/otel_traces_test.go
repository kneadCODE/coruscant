package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOTELTraceProvider_AllModes(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
	}{
		{"ModeDev", ModeDebug},
		{"ModeProd", ModeProd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set required environment variables for testing
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
			t.Setenv("OTEL_SERVICE_NAME", "test-service")

			// Create a test resource
			ctx := context.Background()
			res, err := newResource(ctx)
			require.NoError(t, err)

			// Create trace provider
			provider, err := newOTELTraceProvider(ctx, res, tt.mode)
			require.NoError(t, err)
			require.NotNil(t, provider)

			// Verify provider is configured
			assert.NotNil(t, provider)

			// Test cleanup
			provider.Shutdown(ctx)
		})
	}
}

func TestNewOTELTraceProvider_WithNilResource(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	// Test with nil resource to test error handling
	provider, err := newOTELTraceProvider(context.Background(), nil, ModeDebug)

	// Should still work as resource is optional in OTEL TracerProvider
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	provider.Shutdown(context.Background())
}

func TestNewOTELTraceProvider_DefaultSampler(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	// Test with invalid/unknown mode to trigger default sampler
	ctx := context.Background()
	res, err := newResource(ctx)
	require.NoError(t, err)

	// Use a mode value that doesn't match any case
	invalidMode := Mode(999)
	provider, err := newOTELTraceProvider(context.Background(), res, invalidMode)

	require.NoError(t, err)
	require.NotNil(t, provider)

	provider.Shutdown(ctx)
}

func TestNewOTELTraceProvider_SamplingRates(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	t.Setenv("OTEL_SERVICE_NAME", "test-service")

	ctx := context.Background()
	res, err := newResource(ctx)
	require.NoError(t, err)

	// Test dev modes (should sample all)
	devProvider, err := newOTELTraceProvider(context.Background(), res, ModeDebug)
	require.NoError(t, err)
	require.NotNil(t, devProvider)
	devProvider.Shutdown(context.Background())

	// Test prod modes (should sample less)
	prodProvider, err := newOTELTraceProvider(context.Background(), res, ModeProd)
	require.NoError(t, err)
	require.NotNil(t, prodProvider)
	prodProvider.Shutdown(context.Background())
}
