package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

func TestServerOptionsAndHandlers(t *testing.T) {
	type timeouts struct {
		read, write, idle, graceful time.Duration
	}
	type testCase struct {
		name         string
		options      []ServerOption
		wantErr      bool
		wantAddr     string
		wantTimeouts *timeouts
		testReady    bool
	}
	readinessCalled := false
	readinessHandler := func(w http.ResponseWriter, r *http.Request) {
		readinessCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	}
	cases := []testCase{
		{
			name:     "defaults",
			options:  nil,
			wantErr:  false,
			wantAddr: ":8080",
		},
		{
			name:     "custom port",
			options:  []ServerOption{WithPort(8080)},
			wantErr:  false,
			wantAddr: ":8080",
		},
		{
			name:    "invalid port",
			options: []ServerOption{WithPort(-1)},
			wantErr: true,
		},
		{
			name: "custom timeouts",
			options: []ServerOption{
				WithReadTimeout(1 * time.Second),
				WithWriteTimeout(2 * time.Second),
				WithIdleTimeout(3 * time.Second),
				WithGracefulShutdownTimeout(4 * time.Second),
			},
			wantErr:      false,
			wantTimeouts: &timeouts{1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second},
		},
		{
			name:      "readiness handler",
			options:   []ServerOption{WithReadinessHandler(readinessHandler)},
			wantErr:   false,
			testReady: true,
		},
		{
			name:    "nil readiness handler",
			options: []ServerOption{WithReadinessHandler(nil)},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			readinessCalled = false
			srv, err := NewServer(context.Background(), tc.options...)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, srv)
			if tc.wantAddr != "" {
				assert.Equal(t, tc.wantAddr, srv.srv.Addr)
			}
			if tc.wantTimeouts != nil {
				assert.Equal(t, tc.wantTimeouts.read, srv.srv.ReadTimeout)
				assert.Equal(t, tc.wantTimeouts.write, srv.srv.WriteTimeout)
				assert.Equal(t, tc.wantTimeouts.idle, srv.srv.IdleTimeout)
				assert.Equal(t, tc.wantTimeouts.graceful, srv.gracefulShutdownTimeout)
			}
			if tc.testReady {
				ts := httptest.NewServer(srv.srv.Handler)
				defer ts.Close()
				resp, err := http.Get(ts.URL + "/_/ready")
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.True(t, readinessCalled, "readiness handler should be called")
			}
		})
	}
}

func TestServerStartAndStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	srv, err := NewServer(ctx)
	assert.NoError(t, err, "unexpected error creating server")

	// Start server in goroutine, then cancel context to trigger shutdown
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	err = srv.Start(ctx)
	if err != nil {
		assert.NotContains(t, err.Error(), "http server startup failed: http: Server closed",
			"unexpected error from Start: %v", err)
	}
}

func TestNewServerInvalidOptions(t *testing.T) {
	_, err := NewServer(context.Background(), WithPort(-1))
	assert.Error(t, err, "expected error for invalid port")

	_, err = NewServer(context.Background(), WithReadTimeout(-1))
	assert.Error(t, err, "expected error for invalid read timeout")

	_, err = NewServer(context.Background(), WithWriteTimeout(-1))
	assert.Error(t, err, "expected error for invalid write timeout")

	_, err = NewServer(context.Background(), WithIdleTimeout(-1))
	assert.Error(t, err, "expected error for invalid idle timeout")

	_, err = NewServer(context.Background(), WithGracefulShutdownTimeout(-1))
	assert.Error(t, err, "expected error for invalid graceful shutdown timeout")

	_, err = NewServer(context.Background(), WithReadinessHandler(nil))
	assert.Error(t, err, "expected error for nil readiness handler")

	_, err = NewServer(context.Background(), WithRESTHandler(nil))
	assert.Error(t, err, "expected error for nil REST handler")

	_, err = NewServer(context.Background(), WithGQLHandler(nil))
	assert.Error(t, err, "expected error for nil GQL handler")

	_, err = NewServer(context.Background(), WithMaxHeaderBytes(0))
	assert.Error(t, err, "expected error for invalid max header bytes")
}

func TestCustomHandlersCoverage(t *testing.T) {
	srv, err := NewServer(context.Background(), WithProfilingHandler(), WithRESTHandler(func(r chi.Router) {
		r.Get("/custom", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("custom"))
		})
	}))
	assert.NoError(t, err, "unexpected error creating server")

	ts := httptest.NewServer(srv.srv.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/custom")
	assert.NoError(t, err, "unexpected error making request")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := make([]byte, 6)
	resp.Body.Read(body)
	resp.Body.Close()
	assert.Equal(t, "custom", string(body))
}

func TestPingEndpoint(t *testing.T) {
	srv, err := NewServer(context.Background())
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.srv.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/_/ping")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))

	body := make([]byte, 3)
	resp.Body.Read(body)
	resp.Body.Close()
	assert.Equal(t, "ok\n", string(body))
}

func TestGQLHandlerSuccess(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("gql"))
	})

	srv, err := NewServer(context.Background(), WithGQLHandler(handler))
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.srv.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/graph")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := make([]byte, 3)
	resp.Body.Read(body)
	resp.Body.Close()
	assert.Equal(t, "gql", string(body))
}

func TestMaxHeaderBytesSuccess(t *testing.T) {
	srv, err := NewServer(context.Background(), WithMaxHeaderBytes(2048))
	assert.NoError(t, err)
	assert.Equal(t, 2048, srv.srv.MaxHeaderBytes)
}

func TestNewServerWithNilLogger(t *testing.T) {
	// Test NewServer with context that has no logger
	ctx := context.Background()
	srv, err := NewServer(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, srv)
	assert.Nil(t, srv.srv.ErrorLog)
}

func TestNewServerWithLogger(t *testing.T) {
	// Test NewServer with a context that has a logger to cover the logger path
	ctx := context.Background()

	// Initialize telemetry to get logger in context
	telemetryCtx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDebug)
	if err != nil {
		t.Skip("Could not initialize telemetry for test")
	}
	defer cleanup(ctx)

	srv, err := NewServer(telemetryCtx)
	assert.NoError(t, err)
	assert.NotNil(t, srv)
	assert.NotNil(t, srv.srv.ErrorLog) // Should have a logger now
}

func TestMetricsEndpoint(t *testing.T) {
	srv, err := NewServer(context.Background(), WithMetricsHandler())
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.srv.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/_/metrics")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))

	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	resp.Body.Close()

	content := string(body[:n])
	assert.Contains(t, content, "Metrics are exported via OpenTelemetry")
	assert.Contains(t, content, "See Grafana for comprehensive metrics")
}
