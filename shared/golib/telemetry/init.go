package telemetry

import (
	"context"
	"log"
)

// InitTelemetry initializes telemetry systems and returns a context with telemetry configurations.
func InitTelemetry(ctx context.Context, mode Mode) context.Context {
	log.Println("Initializing telemetry")

	logger := newLogger(mode)
	ctx = setLoggerInContext(ctx, logger)
	logger.DebugContext(ctx, "Telemetry initialized")

	logger.InfoContext(ctx, "Telemetry setup complete")

	return ctx
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
