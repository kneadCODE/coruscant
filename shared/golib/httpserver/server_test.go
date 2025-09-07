package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Start server in goroutine, then cancel context to trigger shutdown
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	err = srv.Start(ctx)
	if err != nil && err.Error() != "http server startup failed: http: Server closed" {
		t.Errorf("unexpected error from Start: %v", err)
	}
}

func TestNewServerInvalidOptions(t *testing.T) {
	_, err := NewServer(context.Background(), WithPort(-1))
	if err == nil {
		t.Error("expected error for invalid port")
	}
	_, err = NewServer(context.Background(), WithReadTimeout(-1))
	if err == nil {
		t.Error("expected error for invalid read timeout")
	}
	_, err = NewServer(context.Background(), WithWriteTimeout(-1))
	if err == nil {
		t.Error("expected error for invalid write timeout")
	}
	_, err = NewServer(context.Background(), WithIdleTimeout(-1))
	if err == nil {
		t.Error("expected error for invalid idle timeout")
	}
	_, err = NewServer(context.Background(), WithGracefulShutdownTimeout(-1))
	if err == nil {
		t.Error("expected error for invalid graceful shutdown timeout")
	}
	_, err = NewServer(context.Background(), WithReadinessHandler(nil))
	if err == nil {
		t.Error("expected error for nil readiness handler")
	}
	_, err = NewServer(context.Background(), WithRESTHandler(nil))
	if err == nil {
		t.Error("expected error for nil REST handler")
	}
	_, err = NewServer(context.Background(), WithGQLHandler(nil))
	if err == nil {
		t.Error("expected error for nil GQL handler")
	}
	_, err = NewServer(context.Background(), WithMaxHeaderBytes(0))
	if err == nil {
		t.Error("expected error for invalid max header bytes")
	}
}

func TestCustomHandlersCoverage(t *testing.T) {
	srv, err := NewServer(context.Background(), WithProfilingHandler(), WithRESTHandler(func(r chi.Router) {
		r.Get("/custom", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("custom"))
		})
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ts := httptest.NewServer(srv.srv.Handler)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/custom")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := make([]byte, 6)
	resp.Body.Read(body)
	resp.Body.Close()
	assert.Equal(t, "custom", string(body))
}
