package telemetry

import (
	"context"
	"log"
	"os"
)

// InitTelemetry initializes telemetry systems and returns a context with telemetry configurations.
// The returned cleanup function should be called during application shutdown.
func InitTelemetry(ctx context.Context, mode Mode) (context.Context, func(), error) {
	log.SetOutput(os.Stdout)
	log.Println("Initializing telemetry")

	// Create OTEL resource once
	resource, err := newResource(ctx)
	if err != nil {
		return ctx, nil, err
	}

	// Initialize logger
	logger, loggerCleanup, err := newLogger(mode, resource)
	if err != nil {
		return ctx, nil, err
	}

	// Initialize trace provider
	traceProvider, traceCleanup, err := newOTELTraceProvider(resource, mode)
	if err != nil {
		loggerCleanup()
		return ctx, nil, err
	}

	// Combined cleanup function
	cleanup := func() {
		traceCleanup()
		loggerCleanup()
	}

	ctx = setLoggerInContext(ctx, logger)
	logger.DebugContext(ctx, "Telemetry initialized", "trace_provider", traceProvider != nil)

	logger.InfoContext(ctx, "Telemetry setup complete")

	return ctx, cleanup, nil
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
