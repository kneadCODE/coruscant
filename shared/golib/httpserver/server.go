package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

// NewServer returns a new instance of Server.
func NewServer(ctx context.Context, options ...ServerOption) (*Server, error) {
	var lg *log.Logger
	if logger := telemetry.LoggerFromContext(ctx); logger != nil {
		lg = slog.NewLogLogger(logger.Handler(), slog.LevelError)
	}
	s := &Server{
		srv: &http.Server{
			Addr:           ":8080",
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1MB default
			ErrorLog:       lg,
			BaseContext: func(net.Listener) context.Context {
				return telemetry.Clone(ctx)
			},
		},
		gracefulShutdownTimeout: 10 * time.Second,
	}

	m := newRouter()

	for _, opt := range options {
		if err := opt(s, m); err != nil {
			return nil, err
		}
	}

	// Add default routes after all middleware is configured
	addDefaultRoutes(m)

	s.srv.Handler = m

	return s, nil
}

// Server is the server instance.
type Server struct {
	srv                     *http.Server
	gracefulShutdownTimeout time.Duration
}

// Start starts the server and is context aware and shuts down when the context gets cancelled.
func (s *Server) Start(ctx context.Context) error {
	startErrChan := make(chan error, 1)

	go func() {
		telemetry.RecordInfoEvent(ctx, fmt.Sprintf("Starting HTTP server on %s", s.srv.Addr))
		startErrChan <- s.srv.ListenAndServe()
	}()

	for {
		select {
		case <-ctx.Done():
			return s.stop(ctx)
		case err := <-startErrChan:
			if !errors.Is(err, http.ErrServerClosed) {
				return fmt.Errorf("http server startup failed: %w", err)
			}
			return nil
		}
	}
}

func (s *Server) stop(ctx context.Context) error {
	cancelCtx, cancel := context.WithTimeout(context.Background(), s.gracefulShutdownTimeout) // Cannot rely on root context as that might have been cancelled.
	defer cancel()

	telemetry.RecordDebugEvent(ctx, "Attempting HTTP server graceful shutdown")
	if err := s.srv.Shutdown(cancelCtx); err != nil {
		telemetry.RecordErrorEvent(ctx, fmt.Errorf("httpserver:Server: graceful shutdown failed: %w", err))

		telemetry.RecordDebugEvent(ctx, "Attempting HTTP server force shutdown")
		if err = s.srv.Close(); err != nil {
			err = fmt.Errorf("httpserver:Server: force shutdown failed: %w", err)
			telemetry.RecordErrorEvent(ctx, err)
			return err
		}
	}

	telemetry.RecordInfoEvent(ctx, "HTTP server shutdown complete")

	return nil
}

// ServerOption customizes the Server.
type ServerOption = func(srv *Server, m *chi.Mux) error

// WithPort overrides the default server port with the given value.
func WithPort(port int) ServerOption {
	return func(s *Server, _ *chi.Mux) error {
		if port <= 0 || port > 65535 {
			return errors.New("invalid port number")
		}
		s.srv.Addr = fmt.Sprintf(":%d", port)
		return nil
	}
}

// WithReadTimeout overrides the default read timeout with the given value.
func WithReadTimeout(d time.Duration) ServerOption {
	return func(s *Server, _ *chi.Mux) error {
		if d < 0 {
			return errors.New("read timeout must be positive")
		}
		s.srv.ReadTimeout = d
		return nil
	}
}

// WithWriteTimeout overrides the default write timeout with the given value.
func WithWriteTimeout(d time.Duration) ServerOption {
	return func(s *Server, _ *chi.Mux) error {
		if d < 0 {
			return errors.New("write timeout must be positive")
		}
		s.srv.WriteTimeout = d
		return nil
	}
}

// WithIdleTimeout overrides the default idle timeout with the given value.
func WithIdleTimeout(d time.Duration) ServerOption {
	return func(s *Server, _ *chi.Mux) error {
		if d < 0 {
			return errors.New("idle timeout must be positive")
		}
		s.srv.IdleTimeout = d
		return nil
	}
}

// WithGracefulShutdownTimeout overrides the default graceful shutdown timeout with the given value.
func WithGracefulShutdownTimeout(d time.Duration) ServerOption {
	return func(s *Server, _ *chi.Mux) error {
		if d < 0 {
			return errors.New("graceful shutdown timeout must be positive")
		}
		s.gracefulShutdownTimeout = d
		return nil
	}
}

// WithProfilingHandler enables go's pprof profiling.
func WithProfilingHandler() ServerOption {
	return func(_ *Server, m *chi.Mux) error {
		// Based on https: //pkg.go.dev/net/http/pprof.
		m.HandleFunc("/_/profile/*", pprof.Index)
		m.HandleFunc("/_/profile/cmdline", pprof.Cmdline)
		m.HandleFunc("/_/profile/profile", pprof.Profile)
		m.HandleFunc("/_/profile/symbol", pprof.Symbol)
		m.HandleFunc("/_/profile/trace", pprof.Trace)
		m.Handle("/_/profile/goroutine", pprof.Handler("goroutine"))
		m.Handle("/_/profile/threadcreate", pprof.Handler("threadcreate"))
		m.Handle("/_/profile/mutex", pprof.Handler("mutex"))
		m.Handle("/_/profile/heap", pprof.Handler("heap"))
		m.Handle("/_/profile/block", pprof.Handler("block"))
		m.Handle("/_/profile/allocs", pprof.Handler("allocs"))
		return nil
	}
}

// WithReadinessHandler sets the handler for readiness checks at `/_/ready`.
func WithReadinessHandler(h http.HandlerFunc) ServerOption {
	return func(_ *Server, m *chi.Mux) error {
		if h == nil {
			return errors.New("readiness handler cannot be nil")
		}
		m.Get("/_/ready", h)
		return nil
	}
}

// WithRESTHandler sets the REST route handler.
func WithRESTHandler(rtr func(chi.Router)) ServerOption {
	return func(_ *Server, m *chi.Mux) error {
		if rtr == nil {
			return errors.New("rest handler cannot be nil")
		}
		m.Group(rtr)
		return nil
	}
}

// WithGQLHandler sets the GQL handler at `/graph` route.
func WithGQLHandler(h http.Handler) ServerOption {
	return func(_ *Server, m *chi.Mux) error {
		if h == nil {
			return errors.New("gql handler cannot be nil")
		}
		m.Handle("/graph", h)
		return nil
	}
}

// WithMaxHeaderBytes sets the maximum allowed size of request headers.
func WithMaxHeaderBytes(n int) ServerOption {
	return func(s *Server, _ *chi.Mux) error {
		if n <= 0 {
			return errors.New("max header bytes must be positive")
		}
		s.srv.MaxHeaderBytes = n
		return nil
	}
}

// newRouter creates and returns a new chi.Mux router without any routes initially.
func newRouter() *chi.Mux {
	rtr := chi.NewRouter()
	rtr.Use(telemetry.HTTPServerTracingMiddleware([]string{"/_/ping", "/_/ready", "/_/health", "/_/metrics"}))
	return rtr
}

// addDefaultRoutes adds the default health endpoints to the router.
func addDefaultRoutes(m *chi.Mux) {
	m.Get("/_/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = fmt.Fprintln(w, "ok") // Intentionally ignoring the error as nothing to do once caught.
	})
}
