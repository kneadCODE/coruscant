package telemetry

import "context"

// Clone returns a new context with the same logger as the original context.
func Clone(ctx context.Context) context.Context {
	newCtx := context.Background()
	newCtx = setLoggerInContext(newCtx, LoggerFromContext(ctx))
	return newCtx
}
