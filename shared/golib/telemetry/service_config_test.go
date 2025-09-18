package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func TestNewServiceConfig(t *testing.T) {
	// Create a resource with common attributes
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test-service"),
			semconv.ServiceVersionKey.String("1.2.3"),
			semconv.ServiceNamespaceKey.String("test-namespace"),
			attribute.String("deployment.environment", "production"),
			semconv.HostNameKey.String("test-host"),
			attribute.String("container.name", "test-container"),
		),
	)
	assert.NoError(t, err)

	config := newServiceConfig(res)

	assert.Equal(t, "test-service", config.Name)
	assert.Equal(t, "test-namespace", config.System)
	assert.Equal(t, "1.2.3", config.Version)
	assert.Equal(t, "production", config.Environment)
	assert.Equal(t, "test-host", config.HostName)
	assert.Equal(t, "test-container", config.ContainerName)
}

func TestNewServiceConfigWithDefaults(t *testing.T) {
	// Create a minimal resource
	res, err := resource.New(context.Background())
	assert.NoError(t, err)

	config := newServiceConfig(res)

	assert.EqualValues(t, ServiceConfig{}, config)
}

func TestServiceConfigContext(t *testing.T) {
	// Create a service config
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("context-test"),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	assert.NoError(t, err)

	config := newServiceConfig(res)

	// Test setting and getting from context
	ctx := context.Background()
	ctx = setServiceConfigInContext(ctx, config)

	retrieved := ServiceConfigFromContext(ctx)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "context-test", retrieved.Name)
	assert.Equal(t, "1.0.0", retrieved.Version)
}

func TestServiceConfigFromContextNil(t *testing.T) {
	// Test with empty context
	ctx := context.Background()
	config := ServiceConfigFromContext(ctx)
	assert.Equal(t, ServiceConfig{}, config)
	assert.False(t, config.IsValid())
}

func TestServiceConfigFromContextWrongType(t *testing.T) {
	// Test with wrong type in context
	ctx := context.WithValue(context.Background(), serviceConfigKey{}, "wrong-type")
	config := ServiceConfigFromContext(ctx)
	assert.Equal(t, ServiceConfig{}, config)
	assert.False(t, config.IsValid())
}

func TestServiceConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   ServiceConfig
		expected bool
	}{
		{
			name:     "valid config with both name and system",
			config:   ServiceConfig{Name: "svc", System: "sys"},
			expected: true,
		},
		{
			name:     "invalid config with empty name",
			config:   ServiceConfig{Name: "", System: "sys"},
			expected: false,
		},
		{
			name:     "invalid config with empty system",
			config:   ServiceConfig{Name: "svc", System: ""},
			expected: false,
		},
		{
			name:     "invalid config with both empty",
			config:   ServiceConfig{Name: "", System: ""},
			expected: false,
		},
		{
			name:     "valid config with additional fields",
			config:   ServiceConfig{Name: "svc", System: "sys", Version: "1.0.0", Environment: "prod"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}
