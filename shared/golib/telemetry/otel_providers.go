package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	olog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// newOTELLogProvider returns a new log provider
func newOTELLogProvider(ctx context.Context, res *resource.Resource) (*olog.LoggerProvider, error) {
	RecordInfoEvent(ctx, "Initializing OTEL gRPC logs exporter")

	// Validate required environment variables
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT environment variable is required")
	}

	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint), // Don't know why but the lib doesn't auto pick up from enevvar
		// No compression for local collector deployment (localhost/same-node)
		// Compression adds CPU overhead without network benefit for local collectors
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTEL gRPC log exporter: %w", err)
	}
	RecordInfoEvent(ctx, "Initialized OTEL gRPC logs exporter")

	// Configure batch processor optimized for local collector deployment
	batchProcessor := olog.NewBatchProcessor(exporter,
		// Balanced settings for local collector (low latency + reasonable throughput)
		olog.WithMaxQueueSize(1024),      // Moderate queue size for local deployment
		olog.WithExportMaxBatchSize(256), // Smaller batches for lower latency to local collector
	)
	logProvider := olog.NewLoggerProvider(
		olog.WithResource(res),
		olog.WithProcessor(batchProcessor),
	)

	RecordInfoEvent(ctx, "Initialized OTEL log provider")

	return logProvider, nil
}

// newOTELTraceProvider creates a new OTEL trace provider with the given resource and mode.
func newOTELTraceProvider(ctx context.Context, res *resource.Resource, mode Mode) (*trace.TracerProvider, error) {
	RecordInfoEvent(ctx, "Initializing OTEL trace gRPC client")

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT environment variable is required")
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint), // Don't know why but the lib doesn't auto pick up from enevvar
		// No compression for local collector deployment (localhost/same-node)
		// Compression adds CPU overhead without network benefit for local collectors)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTEL trace gRPC exporter: %w", err)
	}
	RecordInfoEvent(ctx, "Initialized OTEL trace gRPC exporter")

	// Configure trace provider with appropriate sampling
	var sampler trace.Sampler
	switch mode {
	case ModeDebug:
		// Sample all traces in development
		sampler = trace.AlwaysSample()
	case ModeProd:
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

	RecordInfoEvent(ctx, "Initialized OTEL trace provider")

	return traceProvider, nil
}

// newOTELMetricsProvider creates and configures an OTEL metrics provider.
func newOTELMetricsProvider(ctx context.Context, res *resource.Resource) (*metric.MeterProvider, error) {
	RecordInfoEvent(ctx, "Initializing OTEL gRPC metrics exporter")

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT environment variable is required for metrics")
	}

	// Create OTLP gRPC metrics exporter
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint), // Don't know why but the lib doesn't auto pick up from enevvar
		// No compression for local collector deployment (localhost/same-node)
		// Compression adds CPU overhead without network benefit for local collectors)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTEL gRPC metrics exporter: %w", err)
	}

	RecordInfoEvent(ctx, "Initialized OTEL gRPC metrics exporter")

	// Configure metric reader with optimized settings for local collector deployment
	reader := metric.NewPeriodicReader(exporter,
		// Balanced settings for local collector deployment
		metric.WithInterval(15*time.Second), // 15-second collection interval for local development
		metric.WithTimeout(10*time.Second),  // 10-second timeout for exports
	)

	// Create metrics provider
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)

	RecordInfoEvent(ctx, "Initialized OTEL metrics provider")

	return meterProvider, nil
}
