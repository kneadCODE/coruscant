package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource"
)

// newResource creates an OTEL resource with auto-detected information.
func newResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTEL resource: %w", err)
	}

	return res, nil
}
