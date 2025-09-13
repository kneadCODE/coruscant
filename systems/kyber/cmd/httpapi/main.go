package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

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
	ctx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDevDebug) // TODO: Set the mode as per envvar
	if err != nil {
		return err
	}
	defer cleanup()

	if err := start(ctx); err != nil {
		return err
	}

	return nil
}

func start(ctx context.Context) error {
	// Create HTTP server with tracing middleware enabled for observability
	srv, err := httpserver.NewServer(ctx,
		httpserver.WithTracing(), // Enables distributed tracing for HTTP requests
		httpserver.WithRESTHandler(restHandler),
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
		r.Post("/", testingHandler)
	})
}

func testingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	telemetry.RecordInfoEvent(ctx, "testing endpoint called",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Add some processing time to make the span more visible
	if r.Method == "POST" {
		telemetry.RecordInfoEvent(ctx, "processing POST data")
	}

	w.WriteHeader(200)
	_, _ = w.Write([]byte("OK"))

	telemetry.RecordInfoEvent(ctx, "testing endpoint response sent", "status_code", 200)
}
