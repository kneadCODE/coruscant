package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"

	"github.com/kneadCODE/coruscant/shared/golib/httpserver"
	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	// Initialize telemetry with dev debug mode for comprehensive observability
	ctx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDebug) // TODO: Set the mode as per envvar
	if err != nil {
		return err
	}
	defer cleanup(context.Background())

	if err := start(ctx); err != nil {
		return err
	}

	return nil
}

func start(ctx context.Context) error {
	// Create HTTP server with comprehensive observability features enabled
	srv, err := httpserver.NewServer(ctx,
		httpserver.WithRESTHandler(restHandler),
		httpserver.WithProfilingHandler(), // Enable pprof profiling endpoints
		httpserver.WithMetricsHandler(),   // Enable metrics endpoint (shows OTEL info)
		// Note: HTTP metrics and tracing middleware are automatically enabled in newRouter()
	)
	if err != nil {
		return err
	}
	if err := srv.Start(ctx); err != nil {
		return err
	}
	return nil
}

func restHandler(rtr chi.Router) {
	rtr.Route("/testing", func(r chi.Router) {
		r.Get("/", testingHandler)
		r.Get("/abc", testingHandler)
	})

	rtr.Get("/testing2", testingHandler)
}

func testingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	telemetry.RecordInfoEvent(ctx, "testing endpoint called")

	// Record custom metrics to demonstrate metrics functionality
	if metricsCollector := telemetry.MetricsCollectorFromContext(ctx); metricsCollector != nil {
		// Counter: Track endpoint usage
		metricsCollector.RecordCustomCounter(ctx, "endpoint_calls_total", 1,
			attribute.String("endpoint", "testing"),
			attribute.String("method", r.Method),
		)

		// Gauge: Track active requests (simulated)
		metricsCollector.RecordCustomGauge(ctx, "active_requests", 1,
			attribute.String("endpoint", "testing"),
		)
	}

	someFunc(ctx)

	telemetry.RecordInfoEvent(ctx, "testing endpoint response sent")
}

func someFunc(ctx context.Context) {
	start := time.Now()
	ctx, end := telemetry.Measure(ctx, "response-preparation")
	defer end(nil)

	// Simulate response preparation work
	time.Sleep(10 * time.Millisecond)

	// Record custom histogram for operation duration
	if metricsCollector := telemetry.MetricsCollectorFromContext(ctx); metricsCollector != nil {
		duration := time.Since(start).Seconds()
		metricsCollector.RecordCustomHistogram(ctx, "operation_duration_seconds", duration,
			attribute.String("operation", "response_preparation"),
		)
	}

	telemetry.RecordInfoEvent(ctx, "response prepared")
}
