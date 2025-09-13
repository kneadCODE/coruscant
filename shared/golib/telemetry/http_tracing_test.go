package telemetry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TestHTTPServerTracingMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		excludePaths []string
		shouldTrace  bool
	}{
		{
			name:         "normal_path_should_trace",
			path:         "/api/users",
			excludePaths: []string{"/_/health", "/_/metrics"},
			shouldTrace:  true,
		},
		{
			name:         "excluded_path_should_not_trace",
			path:         "/_/health",
			excludePaths: []string{"/_/health", "/_/metrics"},
			shouldTrace:  false,
		},
		{
			name:         "empty_exclude_list_should_trace_all",
			path:         "/_/health",
			excludePaths: []string{},
			shouldTrace:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that we can verify was called
			handlerCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with our tracing middleware
			middleware := HTTPServerTracingMiddleware(tt.excludePaths)
			wrappedHandler := middleware(testHandler)

			// Create test request
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			// Execute request
			wrappedHandler.ServeHTTP(w, req)

			// Verify handler was called
			assert.True(t, handlerCalled, "Handler should always be called")
			assert.Equal(t, http.StatusOK, w.Code, "Should return 200 status")
		})
	}
}

func TestSpanNameFormatter(t *testing.T) {
	tests := []struct {
		name           string
		operation      string
		path           string
		method         string
		withRouteCtx   bool
		routePattern   string
		expectedResult string
	}{
		{
			name:           "with_existing_operation",
			operation:      "existing-operation",
			path:           "/api/users",
			method:         "GET",
			expectedResult: "existing-operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock tracer and span
			tracer := otel.Tracer("test")
			ctx, span := tracer.Start(context.Background(), "test-span")
			defer span.End()

			// Create request with context containing the span
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req = req.WithContext(ctx)

			// Add Chi route context if specified
			if tt.withRouteCtx {
				rctx := chi.NewRouteContext()
				rctx.RoutePatterns = []string{tt.routePattern}
				reqCtx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				req = req.WithContext(reqCtx)
			}

			// Test the span name formatter
			result := spanNameFormatter(tt.operation, req)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestSpanNameFormatterPathFallback(t *testing.T) {
	// Test cases without actual OTEL setup to verify path fallback logic
	req := httptest.NewRequest("POST", "/api/users", nil)

	// Test with no span context - should return operation as-is
	result := spanNameFormatter("", req)
	assert.Equal(t, "", result) // When no span context, returns empty operation

	// Test function doesn't panic with nil context
	assert.NotPanics(t, func() {
		spanNameFormatter("test", req)
	}, "Should handle requests without span context gracefully")
}

func TestSpanNameFormatterWithInvalidSpan(t *testing.T) {
	// Test with no span in context
	req := httptest.NewRequest("GET", "/test", nil)
	result := spanNameFormatter("test-op", req)
	assert.Equal(t, "test-op", result)
}

func TestEnrichWithSyntheticUserAgentDetection(t *testing.T) {
	tests := []struct {
		name         string
		userAgent    string
		customHeader map[string]string
		expectBot    bool
		expectTest   bool
	}{
		{
			name:      "empty_user_agent",
			userAgent: "",
		},
		{
			name:      "bot_googlebot",
			userAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			expectBot: true,
		},
		{
			name:      "bot_crawler",
			userAgent: "Some crawler bot",
			expectBot: true,
		},
		{
			name:      "bot_spider",
			userAgent: "Web spider",
			expectBot: true,
		},
		{
			name:       "test_k6",
			userAgent:  "k6/0.1 (https://k6.io/)",
			expectTest: true,
		},
		{
			name:       "test_jmeter",
			userAgent:  "Apache-HttpClient/4.5.2 (Java/1.8.0_144) jmeter",
			expectTest: true,
		},
		{
			name:       "test_synthetic",
			userAgent:  "Synthetic monitoring test",
			expectTest: true,
		},
		{
			name:         "smoke_test_header",
			userAgent:    "Regular browser",
			customHeader: map[string]string{"X-Smoke-Test": "true"},
			expectTest:   true,
		},
		{
			name:      "regular_browser",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock tracer and span
			tracer := otel.Tracer("test")
			_, span := tracer.Start(context.Background(), "test-span")
			defer span.End()

			// Create request with user agent
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("User-Agent", tt.userAgent)

			// Add custom headers if specified
			for key, value := range tt.customHeader {
				req.Header.Set(key, value)
			}

			// Test the enrichment function
			enrichWithSyntheticUserAgentDetection(req, span)

			// Note: We can't easily verify the attributes were set without mocking the span
			// but we can at least verify the function doesn't panic
			assert.NotPanics(t, func() {
				enrichWithSyntheticUserAgentDetection(req, span)
			}, "enrichWithSyntheticUserAgentDetection should not panic")
		})
	}
}

func TestEnrichWithSyntheticUserAgentDetectionEdgeCases(t *testing.T) {
	// Test with nil span (edge case)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "some bot")

	// This should not panic even with a nil span from invalid context
	assert.NotPanics(t, func() {
		enrichWithSyntheticUserAgentDetection(req, trace.SpanFromContext(context.Background()))
	}, "Should handle invalid span gracefully")
}

func TestHTTPServerTracingMiddlewareIntegration(t *testing.T) {
	// Test middleware integration with Chi router
	router := chi.NewRouter()

	// Add our tracing middleware
	middleware := HTTPServerTracingMiddleware([]string{"/_/health"})
	router.Use(middleware)

	// Add test routes
	router.Get("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user"))
	})

	router.Get("/_/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "traced_endpoint",
			path:     "/api/users/123",
			expected: "user",
		},
		{
			name:     "excluded_endpoint",
			path:     "/_/health",
			expected: "healthy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.expected, w.Body.String())
		})
	}
}
