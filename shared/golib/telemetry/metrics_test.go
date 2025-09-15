package telemetry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetricsCollector(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_SERVICE_NAME", "test-service")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	require.NoError(t, err)
	defer cleanup(ctx)

	collector, err := NewMetricsCollector()
	require.NoError(t, err)
	assert.NotNil(t, collector)
	assert.NotNil(t, collector.meter)
	assert.NotNil(t, collector.httpServerRequestsInFlight)
	assert.NotNil(t, collector.customCounters)
	assert.NotNil(t, collector.customGauges)
	assert.NotNil(t, collector.customHistograms)
}

func TestHTTPMetricsMiddleware(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_SERVICE_NAME", "test-service")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	require.NoError(t, err)
	if cleanup != nil {
		defer cleanup(ctx)
	}

	collector, err := NewMetricsCollector()
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap handler with metrics middleware
	wrappedHandler := HTTPServerMetricsMiddleware(handler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(setMetricsCollectorInContext(req.Context(), collector))
	rec := httptest.NewRecorder()

	// Make the request
	wrappedHandler.ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestRecordCustomCounter(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_SERVICE_NAME", "test-service")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	require.NoError(t, err)
	if cleanup != nil {
		defer cleanup(ctx)
	}

	collector, err := NewMetricsCollector()
	require.NoError(t, err)

	// Record custom counter metric
	collector.RecordCustomCounter(ctx, "test.requests", 1)
	collector.RecordCustomCounter(ctx, "test.requests", 5, "method", "GET", "status", 200)

	// Should not panic - actual verification would require metric reader
	assert.NotNil(t, collector.customCounters)
}

func TestRecordCustomGauge(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_SERVICE_NAME", "test-service")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	require.NoError(t, err)
	if cleanup != nil {
		defer cleanup(ctx)
	}

	collector, err := NewMetricsCollector()
	require.NoError(t, err)

	// Record custom gauge metric
	collector.RecordCustomGauge(ctx, "test.connections", 10)
	collector.RecordCustomGauge(ctx, "test.connections", 15, "service", "api", "region", "us-west")

	// Should not panic - actual verification would require metric reader
	assert.NotNil(t, collector.customGauges)
}

func TestRecordCustomHistogram(t *testing.T) {
	// Set required environment variables for testing
	t.Setenv("OTEL_SERVICE_NAME", "test-service")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	require.NoError(t, err)
	if cleanup != nil {
		defer cleanup(ctx)
	}

	collector, err := NewMetricsCollector()
	require.NoError(t, err)

	// Record custom histogram metric
	collector.RecordCustomHistogram(ctx, "test.duration", 0.150)
	collector.RecordCustomHistogram(ctx, "test.duration", 0.075, "operation", "read")
	collector.RecordCustomHistogram(ctx, "test.duration", 0.200, "operation", "write", "database", "postgres")

	// Should not panic - actual verification would require metric reader
	assert.NotNil(t, collector.customHistograms)
}

func TestGetStatusClass(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   string
	}{
		{200, "2xx"},
		{201, "2xx"},
		{299, "2xx"},
		{300, "3xx"},
		{301, "3xx"},
		{399, "3xx"},
		{400, "4xx"},
		{404, "4xx"},
		{499, "4xx"},
		{500, "5xx"},
		{503, "5xx"},
		{599, "5xx"},
		{100, "unknown"},
		{600, "5xx"},
	}

	for _, test := range tests {
		result := getStatusClass(test.statusCode)
		assert.Equal(t, test.expected, result, "Status code %d should return %s", test.statusCode, test.expected)
	}
}

func TestGetRoutePattern(t *testing.T) {
	// Test basic path extraction
	req := httptest.NewRequest("GET", "/api/users/123", nil)
	pattern := getRoutePattern(req)
	assert.Equal(t, "/api/users/123", pattern)

	// For full chi router pattern testing, you would need to set up
	// an actual chi router with patterns like "/api/users/{id}"
}

func BenchmarkHTTPMetricsMiddleware(b *testing.B) {
	// Set required environment variables for testing
	b.Setenv("OTEL_SERVICE_NAME", "test-service")
	b.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	require.NoError(b, err)
	if cleanup != nil {
		defer cleanup(ctx)
	}

	collector, err := NewMetricsCollector()
	require.NoError(b, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := HTTPServerMetricsMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(setMetricsCollectorInContext(req.Context(), collector))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
	}
}

func BenchmarkRecordCustomCounter(b *testing.B) {
	// Set required environment variables for testing
	b.Setenv("OTEL_SERVICE_NAME", "test-service")
	b.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=test-system")

	ctx, cleanup, err := InitTelemetry(context.Background(), ModeDebug)
	require.NoError(b, err)
	if cleanup != nil {
		defer cleanup(ctx)
	}

	collector, err := NewMetricsCollector()
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.RecordCustomCounter(ctx, "benchmark.counter", 1, "test", "benchmark")
	}
}

func TestGetServerHost(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		urlHost  string
		expected string
	}{
		{"host header with port", "example.com:8080", "", "example.com"},
		{"host header without port", "example.com", "", "example.com"},
		{"url host with port", "", "api.example.com:9090", "api.example.com"},
		{"url host without port", "", "api.example.com", "api.example.com"},
		{"host header takes precedence", "priority.com", "fallback.com", "priority.com"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://"+test.urlHost+"/test", nil)
			if test.host != "" {
				req.Host = test.host
			}

			result := getServerHost(req)
			assert.Equal(t, test.expected, result)
		})
	}

	// Test empty case separately
	t.Run("empty inputs", func(t *testing.T) {
		req := &http.Request{URL: &url.URL{}}
		result := getServerHost(req)
		assert.Equal(t, "", result)
	})
}

func TestGetServerPort(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		urlHost  string
		scheme   string
		expected int
	}{
		{"host header with port", "example.com:8080", "", "http", 8080},
		{"url host with port", "", "api.example.com:9090", "http", 9090},
		{"https default port", "example.com", "", "https", 443},
		{"http default port", "example.com", "", "http", 80},
		{"host header takes precedence", "priority.com:3000", "fallback.com:4000", "http", 3000},
		{"invalid port", "example.com:invalid", "", "http", 80},
		{"no port info", "example.com", "api.example.com", "http", 80},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", test.scheme+"://"+test.urlHost+"/test", nil)
			if test.host != "" {
				req.Host = test.host
			}

			result := getServerPort(req)
			assert.Equal(t, test.expected, result)
		})
	}
}
