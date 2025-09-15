package telemetry

import (
	"net/http"
	"slices"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/httpconv"
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
			otelhttp.WithMeterProvider(otel.GetMeterProvider()),
			otelhttp.WithMetricAttributesFn(func(r *http.Request) []attribute.KeyValue {
				attrs := []attribute.KeyValue{
					semconv.HTTPRoute(getRoutePattern(r)),
				}
				if attr := getSyntheticUserAgentAttrs(r); attr.Valid() {
					attrs = append(attrs, attr)
				}
				return attrs
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

	attr := getSyntheticUserAgentAttrs(r)
	if attr.Valid() {
		span.SetAttributes(attr)
	}

	// Use existing operation if provided by OTEL
	if operation != "" {
		return operation
	}

	pattern := getRoutePattern(r)

	span.SetAttributes(semconv.HTTPRoute(pattern))

	return r.Method + " " + pattern
}

// getSyntheticUserAgentAttrs detects if a User-Agent string indicates synthetic traffic
func getSyntheticUserAgentAttrs(r *http.Request) attribute.KeyValue {
	ua := strings.ToLower(r.UserAgent())
	if ua == "" {
		return attribute.KeyValue{}
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
			return semconv.UserAgentSyntheticTypeBot
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
			return semconv.UserAgentSyntheticTypeTest
		}
	}

	// Custom smoke test header
	if r.Header.Get("X-Smoke-Test") == "true" {
		return semconv.UserAgentSyntheticTypeTest
	}

	return attribute.KeyValue{}
}

// HTTPServerMetricsMiddleware is a middleware that records HTTP request metrics using OTEL.
func HTTPServerMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c := MetricsCollectorFromContext(ctx)
		if c != nil {
			attrs := []attribute.KeyValue{
				semconv.HTTPRoute(getRoutePattern(r)),
				semconv.URLScheme(r.URL.Scheme),
			}

			// Add server address and port if available
			if host := getServerHost(r); host != "" {
				attrs = append(attrs, semconv.ServerAddress(host))
			}
			if port := getServerPort(r); port > 0 {
				attrs = append(attrs, semconv.ServerPort(port))
			}

			// Track in-flight requests
			c.httpServerRequestsInFlight.Add(ctx, 1,
				httpconv.RequestMethodAttr(r.Method),
				r.URL.Scheme,
				attrs...,
			)
			// Decrement in-flight requests
			defer func() {
				c.httpServerRequestsInFlight.Add(ctx, -1,
					httpconv.RequestMethodAttr(r.Method),
					r.URL.Scheme,
					attrs...,
				)
			}()
		}

		// Process request
		next.ServeHTTP(w, r)
	})
}
