package telemetry

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// newOTELTraceProvider creates a new OTEL trace provider with the given resource and mode.
func newOTELTraceProvider(ctx context.Context, res *resource.Resource, mode Mode) (*trace.TracerProvider, func(context.Context), error) {
	RecordInfoEvent(ctx, "Initializing OTEL trace gRPC client")

	// Validate required environment variables
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil, nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT environment variable is required")
	}

	// Configure OTLP options based on environment
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
		// No compression for local collector deployment (localhost/same-node)
		// Compression adds CPU overhead without network benefit for local collectors
	}

	// Check if insecure connection is requested (for local development)
	if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true" {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OTEL trace gRPC exporter: %w", err)
	}
	RecordInfoEvent(ctx, "Initialized OTEL trace gRPC exporter")

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

	// Configure trace provider optimized for local collector deployment
	traceProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exporter,
			// Balanced settings for local collector (low latency + reasonable throughput)
			trace.WithMaxQueueSize(1024),      // Moderate queue size for local deployment
			trace.WithMaxExportBatchSize(256), // Smaller batches for lower latency to local collector
		),
		trace.WithSampler(sampler),
	)

	// Set global trace provider
	otel.SetTracerProvider(traceProvider)

	RecordInfoEvent(ctx, "Initialized OTEL trace provider")

	return traceProvider, func(ctx context.Context) {
		_ = traceProvider.Shutdown(ctx)
	}, nil
}
