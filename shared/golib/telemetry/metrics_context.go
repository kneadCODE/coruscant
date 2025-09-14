package telemetry

import (
	"context"
)

// metricsCollectorKey is the context key for the metrics collector.
type metricsCollectorKey struct{}

// setMetricsCollectorInContext sets the metrics collector in the context.
func setMetricsCollectorInContext(ctx context.Context, collector *MetricsCollector) context.Context {
	return context.WithValue(ctx, metricsCollectorKey{}, collector)
}

// MetricsCollectorFromContext retrieves the metrics collector from the context.
// Returns nil if no collector is found in the context.
func MetricsCollectorFromContext(ctx context.Context) *MetricsCollector {
	collector, ok := ctx.Value(metricsCollectorKey{}).(*MetricsCollector)
	if !ok {
		return nil
	}
	return collector
}
