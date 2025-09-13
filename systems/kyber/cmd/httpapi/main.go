package main

import (
	"context"
	"log"
	"net/http"
	"time"

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
	defer cleanup(context.Background())

	if err := start(ctx); err != nil {
		return err
	}

	return nil
}

func start(ctx context.Context) error {
	// Create HTTP server with tracing middleware enabled for observability
	srv, err := httpserver.NewServer(ctx,
		// httpserver.WithTracing(),
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
		r.Get("/abc", testingHandler)
	})

	rtr.Get("/testing2", testingHandler)
}

func testingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	telemetry.RecordInfoEvent(ctx, "testing endpoint called")

	someFunc(ctx)

	telemetry.RecordInfoEvent(ctx, "testing endpoint response sent")
}

func someFunc(ctx context.Context) {
	ctx, end := telemetry.Measure(ctx, "response-preparation")
	defer end(nil)

	// Simulate response preparation work
	time.Sleep(10 * time.Millisecond)

	telemetry.RecordInfoEvent(ctx, "response prepared")
}
