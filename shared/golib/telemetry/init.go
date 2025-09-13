package telemetry

import (
	"context"
	"log"
	"os"
)

// InitTelemetry initializes telemetry systems and returns a context with telemetry configurations.
// The returned cleanup function should be called during application shutdown.
func InitTelemetry(ctx context.Context, mode Mode) (context.Context, func(context.Context), error) {
	log.SetOutput(os.Stdout)
	log.Println("Initializing telemetry")

	cleanupF := func(context.Context) {}

	// Create OTEL resource once
	resource, err := newResource(ctx)
	if err != nil {
		return ctx, cleanupF, err
	}

	// Initialize logger
	logger, loggerCleanup, err := newOTELSlogLogger(ctx, resource)
	if err != nil {
		return ctx, cleanupF, err
	}
	ctx = setLoggerInContext(ctx, logger)
	cleanupF = func(ctx context.Context) {
		loggerCleanup(ctx)
	}

	// Initialize trace provider
	_, traceCleanup, err := newOTELTraceProvider(ctx, resource, mode)
	if err != nil {
		return ctx, cleanupF, err
	}

	// Combined cleanup function
	cleanupF = func(ctx context.Context) {
		traceCleanup(ctx)
		loggerCleanup(ctx)
	}

	logger.InfoContext(ctx, "Telemetry initialization complete")

	return ctx, cleanupF, nil
}

// Mode represents the telemetry/logging mode.
type Mode int

const (
	// ModeDev enables development logging (info level, text output).
	ModeDev Mode = iota
	// ModeDevDebug enables development logging with debug level.
	ModeDevDebug
	// ModeProd enables production logging (info level, JSON output).
	ModeProd
	// ModeProdDebug enables production logging with debug level.
	ModeProdDebug
)
