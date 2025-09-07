// Package telemetry enables context-aware structured logging for Go using slog.
//
//	// Record events
//	telemetry.RecordInfoEvent(ctx, "request started", "user_id", userID)
//	telemetry.RecordErrorEvent(ctx, err, "user_id", userID)
//	telemetry.RecordDebugEvent(ctx, "debug details", "foo", "bar")
//
//	// Use slog's context methods
//	slog.InfoContext(ctx, "message") // uses logger from context
package telemetry
