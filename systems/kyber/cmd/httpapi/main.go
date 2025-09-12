package main

import (
	"context"
	"log"

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
	srv, err := httpserver.NewServer(ctx)
	if err != nil {
		return err
	}
	if err := srv.Start(ctx); err != nil {
		return err
	}
	return nil
}
