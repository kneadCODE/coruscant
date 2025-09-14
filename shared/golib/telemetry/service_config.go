package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// ServiceConfig holds frequently used service information extracted from OTEL resource attributes.
// This avoids repeatedly iterating over resource attributes throughout the application.
type ServiceConfig struct {
	// Core service identity
	Name    string
	System  string
	Version string

	// Deployment information
	Environment string

	// Infrastructure details
	HostName      string
	ContainerName string

	// Additional common attributes can be added as needed
}

// IsValid checks if the essential fields of ServiceConfig are populated.
func (sc ServiceConfig) IsValid() bool {
	return sc.Name != "" && sc.System != ""
}

// serviceConfigKey is the context key for the service configuration.
type serviceConfigKey struct{}

// newServiceConfig creates a ServiceConfig by extracting common attributes from the OTEL resource.
func newServiceConfig(res *resource.Resource) ServiceConfig {
	config := ServiceConfig{}

	// Extract common attributes from resource
	for _, kv := range res.Attributes() {
		switch kv.Key {
		case semconv.ServiceNameKey:
			config.Name = kv.Value.AsString()
		case semconv.ServiceNamespaceKey:
			config.System = kv.Value.AsString()
		case semconv.ServiceVersionKey:
			config.Version = kv.Value.AsString()
		case semconv.DeploymentEnvironmentNameKey:
			config.Environment = kv.Value.AsString()
		case semconv.HostNameKey:
			config.HostName = kv.Value.AsString()
		case semconv.ContainerNameKey:
			config.ContainerName = kv.Value.AsString()
		}

		// Handle custom deployment.environment attribute for backward compatibility
		if string(kv.Key) == "deployment.environment" {
			config.Environment = kv.Value.AsString()
		}
		// Handle custom container.name attribute for backward compatibility
		if string(kv.Key) == "container.name" {
			config.ContainerName = kv.Value.AsString()
		}
	}

	return config
}

// setServiceConfigInContext stores the service configuration in the context.
func setServiceConfigInContext(ctx context.Context, config ServiceConfig) context.Context {
	return context.WithValue(ctx, serviceConfigKey{}, config)
}

// ServiceConfigFromContext retrieves the service configuration from the context.
// Returns zero value if no configuration is found in the context.
func ServiceConfigFromContext(ctx context.Context) ServiceConfig {
	config, ok := ctx.Value(serviceConfigKey{}).(ServiceConfig)
	if !ok {
		return ServiceConfig{}
	}
	return config
}
