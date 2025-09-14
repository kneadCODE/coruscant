package telemetry

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv/v1.37.0/httpconv"
)

// MetricsCollector provides OTEL-based metrics collection integrated with telemetry.
type MetricsCollector struct {
	meter                      metric.Meter                         // OTEL meter instance for metric creation
	httpServerRequestsInFlight httpconv.ServerActiveRequests        // Tracks active HTTP requests
	customCounters             map[string]metric.Int64Counter       // Dynamic counter metrics registry
	customGauges               map[string]metric.Int64UpDownCounter // Dynamic gauge metrics registry
	customHistograms           map[string]metric.Float64Histogram   // Dynamic histogram metrics registry
	metricsMutex               sync.RWMutex                         // Protects concurrent access to metric maps
}

// NewMetricsCollector creates a new metrics collector using OTEL Metrics API.
// Returns nil and error if initialization fails.
func NewMetricsCollector() (*MetricsCollector, error) {
	meter := otel.Meter(instrumentationIdentifier)

	c := &MetricsCollector{
		meter:            meter,
		customCounters:   make(map[string]metric.Int64Counter),
		customGauges:     make(map[string]metric.Int64UpDownCounter),
		customHistograms: make(map[string]metric.Float64Histogram),
	}

	if err := c.initHTTPMetrics(); err != nil {
		return nil, err
	}
	if err := c.initRuntimeMetrics(); err != nil {
		return nil, err
	}

	return c, nil
}

// initHTTPMetrics initializes HTTP-related metrics using OTEL instruments.
func (c *MetricsCollector) initHTTPMetrics() error {
	var err error

	// In-flight requests gauge
	c.httpServerRequestsInFlight, err = httpconv.NewServerActiveRequests(c.meter)
	if err != nil {
		return fmt.Errorf("failed to create http.server.active_requests metric: %w", err)
	}
	return nil
}

// initRuntimeMetrics initializes Go runtime metrics using the contrib runtime package.
func (c *MetricsCollector) initRuntimeMetrics() error {
	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(15 * time.Second))
	if err != nil {
		return fmt.Errorf("failed to start runtime metrics: %w", err)
	}
	return nil
}

// RecordCustomCounter records a custom counter metric with OTEL attributes.
func (c *MetricsCollector) RecordCustomCounter(ctx context.Context, name string, value int64, attrs ...attribute.KeyValue) {
	c.metricsMutex.Lock()
	defer c.metricsMutex.Unlock()

	counter, exists := c.customCounters[name]
	if !exists {
		var err error
		counter, err = c.meter.Int64Counter(
			"custom."+name,
			metric.WithDescription("Custom counter metric: "+name),
			metric.WithUnit("1"),
		)
		if err != nil {
			return // Silently fail to avoid disrupting application
		}
		c.customCounters[name] = counter
	}

	counter.Add(ctx, value, metric.WithAttributes(attrs...))
}

// RecordCustomGauge records a custom gauge metric with OTEL attributes.
func (c *MetricsCollector) RecordCustomGauge(ctx context.Context, name string, value int64, attrs ...attribute.KeyValue) {
	c.metricsMutex.Lock()
	defer c.metricsMutex.Unlock()

	gauge, exists := c.customGauges[name]
	if !exists {
		var err error
		gauge, err = c.meter.Int64UpDownCounter(
			"custom."+name,
			metric.WithDescription("Custom gauge metric: "+name),
			metric.WithUnit("1"),
		)
		if err != nil {
			return // Silently fail to avoid disrupting application
		}
		c.customGauges[name] = gauge
	}

	// Note: OTEL UpDownCounter doesn't have a Set method
	// For custom gauges, we'll use Add with the full value each time
	// In a production system, you'd want to track previous values and calculate deltas
	gauge.Add(ctx, value, metric.WithAttributes(attrs...))
}

// RecordCustomHistogram records a custom histogram metric with OTEL attributes.
func (c *MetricsCollector) RecordCustomHistogram(ctx context.Context, name string, value float64, attrs ...attribute.KeyValue) {
	c.metricsMutex.Lock()
	defer c.metricsMutex.Unlock()

	histogram, exists := c.customHistograms[name]
	if !exists {
		var err error
		histogram, err = c.meter.Float64Histogram(
			"custom."+name,
			metric.WithDescription("Custom histogram metric: "+name),
		)
		if err != nil {
			return // Silently fail to avoid disrupting application
		}
		c.customHistograms[name] = histogram
	}

	histogram.Record(ctx, value, metric.WithAttributes(attrs...))
}

// getStatusClass returns HTTP status class (2xx, 3xx, 4xx, 5xx).
func getStatusClass(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "2xx"
	case statusCode >= 300 && statusCode < 400:
		return "3xx"
	case statusCode >= 400 && statusCode < 500:
		return "4xx"
	case statusCode >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}

// getRoutePattern extracts route pattern for consistent metrics labeling.
func getRoutePattern(r *http.Request) string {
	// Try to get the Chi route pattern for better span grouping
	if rctx := chi.RouteContext(r.Context()); rctx != nil && rctx.RoutePattern() != "" {
		// Use route pattern (e.g., "/api/users/{id}") for better aggregation
		return rctx.RoutePattern()
	}

	// Fallback to actual path if no route pattern available
	return r.URL.Path
}

// getServerHost extracts the server host from the request.
func getServerHost(r *http.Request) string {
	// Prefer the Host header
	if r.Host != "" {
		// Split host:port if needed
		if host, _, err := net.SplitHostPort(r.Host); err == nil {
			return host
		}
		return r.Host
	}

	// Fallback to URL host
	if r.URL.Host != "" {
		if host, _, err := net.SplitHostPort(r.URL.Host); err == nil {
			return host
		}
		return r.URL.Host
	}

	return ""
}

// getServerPort extracts the server port from the request.
func getServerPort(r *http.Request) int {
	// Try to get port from Host header
	if r.Host != "" {
		if _, portStr, err := net.SplitHostPort(r.Host); err == nil {
			if port, err := strconv.Atoi(portStr); err == nil {
				return port
			}
		}
	}

	// Try to get port from URL
	if r.URL.Host != "" {
		if _, portStr, err := net.SplitHostPort(r.URL.Host); err == nil {
			if port, err := strconv.Atoi(portStr); err == nil {
				return port
			}
		}
	}

	// Infer default ports based on scheme
	switch strings.ToLower(r.URL.Scheme) {
	case "https":
		return 443
	case "http":
		return 80
	}

	return 0
}
