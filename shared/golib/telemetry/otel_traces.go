package telemetry

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// newOTELTraceProvider creates a new OTEL trace provider with the given resource and mode.
func newOTELTraceProvider(res *resource.Resource, mode Mode) (*trace.TracerProvider, func(), error) {
	// For prod modes, use stdout exporter (in real production, this would be OTLP)
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(os.Stdout),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
	}

	// Configure trace provider with appropriate sampling
	var sampler trace.Sampler
	switch mode {
	case ModeDev, ModeDevDebug:
		// Sample all traces in development
		sampler = trace.AlwaysSample()
	case ModeProd, ModeProdDebug:
		// Sample 10% of traces in production (could enhance with smart sampling)
		sampler = trace.ParentBased(trace.TraceIDRatioBased(0.1))
	default:
		sampler = trace.AlwaysSample()
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exporter),
		trace.WithSampler(sampler),
	)

	// Set global trace provider
	otel.SetTracerProvider(traceProvider)

	cleanup := func() {
		_ = traceProvider.Shutdown(context.Background())
	}

	return traceProvider, cleanup, nil
}
