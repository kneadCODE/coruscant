package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

// newOTELSlogHandler creates a new slog handler that sends logs to OTEL.
func newOTELSlogHandler(res *resource.Resource, mode Mode) (slog.Handler, func(), error) {
	// For prod modes, use no exporter (will need OTLP exporter later)
	exporter, err := stdoutlog.New(
		stdoutlog.WithWriter(os.Stdout),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed stdoutlog exporter init: %w", err)
	}

	logProvider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(log.NewBatchProcessor(exporter)),
	)

	// Create the bridge handler
	handler := otelslog.NewHandler("service", otelslog.WithLoggerProvider(logProvider))

	return handler, func() {
		_ = logProvider.Shutdown(context.Background())
	}, nil
}
