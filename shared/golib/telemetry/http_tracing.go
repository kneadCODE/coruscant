package telemetry

import (
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

// HTTPServerTracingMiddleware returns an HTTP middleware that provides enhanced OpenTelemetry tracing
// with comprehensive HTTP semantic conventions including missing attributes like route patterns,
// request/response sizes, and additional network information.
func HTTPServerTracingMiddleware(excludePaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(
			next,
			"",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			otelhttp.WithSpanNameFormatter(spanNameFormatter),
			otelhttp.WithFilter(func(r *http.Request) bool {
				path := r.URL.Path
				return !slices.Contains(excludePaths, path)
			}),
		)
	}
}

// Add enhanced span name formatter with Chi route pattern support and enrich with missing attributes
func spanNameFormatter(operation string, r *http.Request) string {
	span := trace.SpanFromContext(r.Context())
	if span == nil || !span.SpanContext().IsValid() {
		return operation
	}

	enrichWithSyntheticUserAgentDetection(r, span)

	// Use existing operation if provided by OTEL
	if operation != "" {
		return operation
	}

	// Try to get the Chi route pattern for better span grouping
	if rctx := chi.RouteContext(r.Context()); rctx != nil && rctx.RoutePattern() != "" {
		// Use route pattern (e.g., "/api/users/{id}") for better aggregation
		operation = r.Method + " " + rctx.RoutePattern()
		span.SetAttributes(semconv.HTTPRoute(rctx.RoutePattern()))
	} else {
		// Fallback to actual path if no route pattern available
		operation = r.Method + " " + r.URL.Path
	}

	return operation
}

// enrichWithSyntheticUserAgentDetection detects if a User-Agent string indicates synthetic traffic
func enrichWithSyntheticUserAgentDetection(r *http.Request, span trace.Span) {
	ua := strings.ToLower(r.UserAgent())
	if ua == "" {
		return
	}

	// Bot indicators (crawlers, scrapers, validators)
	botIndicators := []string{
		"bot", "crawler", "spider", "scraper", "scanner",
		"validator", "headless", "phantom", "selenium", "webdriver",
		"googlebot", "bingbot", "slurp", "duckduckbot",
		"facebookexternalhit", "twitterbot", "linkedinbot",
	}
	for _, b := range botIndicators {
		if strings.Contains(ua, b) {
			span.SetAttributes(semconv.UserAgentSyntheticTypeBot)
			return
		}
	}

	// Test indicators (synthetic monitoring, load testing, probes)
	testIndicators := []string{
		"k6", "jmeter", "locust", "wrk", "ab", "synthetic",
		"pingdom", "grafanacloud/syntheticmonitoring",
		"googlehc", "healthchecker",
	}
	for _, t := range testIndicators {
		if strings.Contains(ua, t) {
			span.SetAttributes(semconv.UserAgentSyntheticTypeTest)
			return
		}
	}

	// Custom smoke test header
	if r.Header.Get("X-Smoke-Test") == "true" {
		span.SetAttributes(semconv.UserAgentSyntheticTypeTest)
	}
}
