package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// TracingConfig configures HTTP tracing behavior.
type TracingConfig struct {
	// IncludeMessageEvents controls whether to trace HTTP request/response body read/write events.
	// When true, adds detailed I/O events to spans (useful for debugging, but more verbose).
	// When false, only traces the overall HTTP request without body I/O details.
	IncludeMessageEvents bool
	// AdditionalFilteredPaths specifies additional paths to exclude from tracing.
	// Health endpoints (/_/ping, /_/ready, /_/health, /_/metrics) are always filtered.
	AdditionalFilteredPaths []string
}

// WithTracing enables OpenTelemetry HTTP tracing for all requests with default configuration.
// This should be enabled to get distributed tracing across HTTP requests.
// Health endpoints are automatically filtered from tracing.
func WithTracing() ServerOption {
	return WithTracingConfig(TracingConfig{
		IncludeMessageEvents: true,
	})
}

// WithTracingConfig enables OpenTelemetry HTTP tracing with custom configuration.
func WithTracingConfig(config TracingConfig) ServerOption {
	return func(_ *Server, m *chi.Mux) error {
		// Add OTEL HTTP middleware to the router
		m.Use(func(next http.Handler) http.Handler {
			var opts []otelhttp.Option

			// Add message events if requested
			if config.IncludeMessageEvents {
				opts = append(opts, otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))
			}

			// Add span name formatter to use route patterns instead of actual paths
			opts = append(opts, otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				// Use existing operation if provided by OTEL
				if operation != "" {
					return operation
				}

				// Try to get the Chi route pattern for better span grouping
				if rctx := chi.RouteContext(r.Context()); rctx != nil && rctx.RoutePattern() != "" {
					// Use route pattern (e.g., "/api/users/{id}") for better aggregation
					return r.Method + " " + rctx.RoutePattern()
				}

				// Fallback to actual path if no route pattern available
				return r.Method + " " + r.URL.Path
			}))

			// Always filter health endpoints + any additional paths
			opts = append(opts, otelhttp.WithFilter(func(r *http.Request) bool {
				path := r.URL.Path

				// Always filter health/monitoring endpoints - never trace these
				if path == "/_/ping" || path == "/_/ready" || path == "/_/health" || path == "/_/metrics" {
					return false
				}

				// Filter additional user-specified paths
				for _, filteredPath := range config.AdditionalFilteredPaths {
					if path == filteredPath {
						return false
					}
				}

				return true
			}))

			return otelhttp.NewHandler(next, "", opts...)
		})
		return nil
	}
}
