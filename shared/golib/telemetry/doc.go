// Package telemetry provides OpenTelemetry-based observability with structured logging and distributed tracing.
//
// It integrates slog for structured logging and OpenTelemetry for traces, all context-aware.
// Events are logged and also recorded as span events when tracing is active.
//
// Basic usage:
//
//	// Initialize telemetry (typically in main)
//	ctx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDevDebug)
//	if err != nil {
//		return err
//	}
//	defer cleanup()
//
//	// Record events (logged and added to spans)
//	telemetry.RecordInfoEvent(ctx, "request started", "user_id", userID)
//	telemetry.RecordErrorEvent(ctx, err, "user_id", userID)
//	telemetry.RecordDebugEvent(ctx, "debug details", "foo", "bar")
//
//	// Use slog's context methods directly
//	if logger := telemetry.LoggerFromContext(ctx); logger != nil {
//		logger.InfoContext(ctx, "message") // uses logger from context
//	}
//
// Modes:
//   - ModeDev/ModeDevDebug: Samples all traces, suitable for development
//   - ModeProd/ModeProdDebug: Samples 10% of traces, suitable for production
package telemetry
