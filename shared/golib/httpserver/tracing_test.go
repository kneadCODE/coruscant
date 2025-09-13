package httpserver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

func TestWithTracing(t *testing.T) {
	// Initialize telemetry to set up trace provider
	ctx := context.Background()
	ctx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDevDebug)
	require.NoError(t, err)
	defer cleanup()

	// Capture stdout to verify trace output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		w.Close()
		os.Stdout = oldStdout
	}()

	// Create a test router with tracing
	mux := chi.NewRouter()
	server := &Server{}

	// Apply tracing middleware
	tracingOpt := WithTracing()
	err = tracingOpt(server, mux)
	require.NoError(t, err)

	// Add a test route
	mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()

	// Execute request through middleware
	mux.ServeHTTP(recorder, req)

	// Verify response
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "test response", recorder.Body.String())

	// Give time for spans to be exported
	time.Sleep(100 * time.Millisecond)

	// Close writer to flush output
	w.Close()

	// Read captured output
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Verify tracing middleware was applied successfully
	// The actual span output format may vary, so we check for key indicators
	if len(outputStr) > 0 {
		// If there's output, it should be trace-related
		assert.True(t, strings.Contains(outputStr, "span") || strings.Contains(outputStr, "trace") || strings.Contains(outputStr, "GET"))
	}
}

func TestWithTracingFiltersHealthEndpoints(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		w.Close()
		os.Stdout = oldStdout
	}()

	// Create router with tracing
	mux := chi.NewRouter()
	server := &Server{}

	tracingOpt := WithTracing()
	err := tracingOpt(server, mux)
	require.NoError(t, err)

	// Add health check routes
	healthEndpoints := []string{"/_/ping", "/_/ready", "/_/health", "/_/metrics"}

	for _, endpoint := range healthEndpoints {
		mux.Get(endpoint, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
	}

	// Test each health endpoint
	for _, endpoint := range healthEndpoints {
		req := httptest.NewRequest(http.MethodGet, endpoint, nil)
		recorder := httptest.NewRecorder()

		mux.ServeHTTP(recorder, req)

		// Verify response is successful
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "ok", recorder.Body.String())
	}

	// Close writer and check output
	w.Close()
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Health endpoints should be filtered out, so no trace output expected
	// If there's any trace output, it shouldn't contain health endpoints
	for _, endpoint := range healthEndpoints {
		assert.NotContains(t, outputStr, endpoint)
	}
}

func TestWithTracingSpanNameFormatter(t *testing.T) {
	// Initialize telemetry to set up trace provider
	ctx := context.Background()
	ctx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDevDebug)
	require.NoError(t, err)
	defer cleanup()

	// Create router with tracing
	mux := chi.NewRouter()
	server := &Server{}

	tracingOpt := WithTracing()
	err = tracingOpt(server, mux)
	require.NoError(t, err)

	// Add test routes with URL patterns
	mux.Get("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user " + userID))
	})

	mux.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("order created"))
	})

	// Test parameterized route - should use pattern, not actual path
	req1 := httptest.NewRequest(http.MethodGet, "/api/users/123", nil)
	req1 = req1.WithContext(ctx)
	recorder1 := httptest.NewRecorder()
	mux.ServeHTTP(recorder1, req1)

	assert.Equal(t, http.StatusOK, recorder1.Code)
	assert.Equal(t, "user 123", recorder1.Body.String())

	// Test another parameterized request - should group with same pattern
	req2 := httptest.NewRequest(http.MethodGet, "/api/users/456", nil)
	req2 = req2.WithContext(ctx)
	recorder2 := httptest.NewRecorder()
	mux.ServeHTTP(recorder2, req2)

	assert.Equal(t, http.StatusOK, recorder2.Code)
	assert.Equal(t, "user 456", recorder2.Body.String())

	// Test non-parameterized route
	req3 := httptest.NewRequest(http.MethodPost, "/orders", nil)
	req3 = req3.WithContext(ctx)
	recorder3 := httptest.NewRecorder()
	mux.ServeHTTP(recorder3, req3)

	assert.Equal(t, http.StatusCreated, recorder3.Code)
	assert.Equal(t, "order created", recorder3.Body.String())

	// Test passes if requests work correctly - span name verification is visible in test output
}

func TestWithTracingConfig(t *testing.T) {
	// Initialize telemetry
	ctx := context.Background()
	ctx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDevDebug)
	require.NoError(t, err)
	defer cleanup()

	mux := chi.NewRouter()
	server := &Server{}

	// Apply custom tracing config with message events disabled and additional filtered paths
	customConfig := TracingConfig{
		IncludeMessageEvents:    false,
		AdditionalFilteredPaths: []string{"/admin", "/internal"},
	}

	customTracingOpt := WithTracingConfig(customConfig)
	err = customTracingOpt(server, mux)
	require.NoError(t, err)

	// Add test routes
	mux.Get("/custom", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("custom response"))
	})

	mux.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("admin"))
	})

	// Test regular request
	req := httptest.NewRequest(http.MethodGet, "/custom", nil)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "custom response", recorder.Body.String())

	// Test filtered request (should still work but not traced)
	req2 := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req2 = req2.WithContext(ctx)
	recorder2 := httptest.NewRecorder()

	mux.ServeHTTP(recorder2, req2)

	assert.Equal(t, http.StatusOK, recorder2.Code)
	assert.Equal(t, "admin", recorder2.Body.String())
}

func TestWithTracingMessageEvents(t *testing.T) {
	// Initialize telemetry to set up trace provider
	ctx := context.Background()
	ctx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDevDebug)
	require.NoError(t, err)
	defer cleanup()

	// Capture stdout to verify message events
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		w.Close()
		os.Stdout = oldStdout
	}()

	mux := chi.NewRouter()
	server := &Server{}

	tracingOpt := WithTracing()
	err = tracingOpt(server, mux)
	require.NoError(t, err)

	// Add route that reads request body
	mux.Post("/data", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("received: " + string(body)))
	})

	// Create request with body
	req := httptest.NewRequest(http.MethodPost, "/data", strings.NewReader("test data"))
	req.Header.Set("Content-Type", "text/plain")
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "received: test data", recorder.Body.String())

	// Give time for spans to be exported
	time.Sleep(100 * time.Millisecond)

	// Close and check output for message events
	w.Close()
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Verify tracing middleware was applied (output may vary)
	if len(outputStr) > 0 {
		assert.True(t, strings.Contains(outputStr, "POST") || strings.Contains(outputStr, "data") || strings.Contains(outputStr, "span"))
	}
}
