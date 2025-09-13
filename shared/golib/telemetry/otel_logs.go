package telemetry

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	olog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

// newOTELSlogHandler creates a new slog logger that sends logs to OTEL.
func newOTELSlogLogger(ctx context.Context, res *resource.Resource) (*slog.Logger, func(context.Context), error) {
	var exporter olog.Exporter
	var err error

	// Validate required environment variables
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil, nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT environment variable is required")
	}

	log.Println("Initializing OTEL gRPC logs exporter")

	// Configure OTLP options based on environment
	opts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(endpoint),
		// No compression for local collector deployment (localhost/same-node)
		// Compression adds CPU overhead without network benefit for local collectors
	}

	// Check if insecure connection is requested (for local development)
	if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true" {
		opts = append(opts, otlploggrpc.WithInsecure())
	}

	exporter, err = otlploggrpc.New(ctx, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OTEL gRPC log exporter: %w", err)
	}
	log.Println("Initialized OTEL gRPC logs exporter")

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

	handler := otelslog.NewHandler(
		instrumentationIdentifier,
		otelslog.WithLoggerProvider(logProvider),
	)

	logger := slog.New(handler)

	logger.InfoContext(ctx, "Initialized OTEL slog logger")

	return logger, func(ctx context.Context) {
		_ = logProvider.Shutdown(ctx)
	}, nil
}
