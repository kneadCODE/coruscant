package main

import (
	"context"
	"log"

	"github.com/kneadCODE/coruscant/shared/golib/httpserver"
	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

func main() {
	run(context.Background())
}

func run(ctx context.Context) {
	ctx = telemetry.InitTelemetry(ctx, telemetry.ModeDevDebug) // TODO: Set the mode as per envvar

	if err := start(ctx); err != nil {
		log.Fatal(err)
	}
}

func start(ctx context.Context) error {
	srv, err := httpserver.NewServer(ctx)
	if err != nil {
		return err
	}
	if err := srv.Start(ctx); err != nil {
		return err
	}
	return nil
}
