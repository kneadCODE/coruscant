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
func newOTELSlogLogger(ctx context.Context, res *resource.Resource) (*slog.Logger, func(), error) {
	var exporter olog.Exporter
	var err error

	log.Println("Initializing OTEL gRPC logs exporter")
	exporter, err = otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
		// No compression for local collector deployment (localhost/same-node)
		// Compression adds CPU overhead without network benefit for local collectors
	)
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
		os.Getenv("OTEL_SERVICE_NAME"),
		otelslog.WithLoggerProvider(logProvider),
	)

	logger := slog.New(handler)

	logger.InfoContext(ctx, "Initialized OTEL slog logger")

	return logger, func() {
		_ = logProvider.Shutdown(context.Background())
	}, nil
}
