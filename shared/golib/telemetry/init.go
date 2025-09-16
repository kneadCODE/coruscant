package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	otelpyroscope "github.com/grafana/otel-profiling-go"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
)

// InitTelemetry initializes telemetry systems and returns a context with telemetry configurations.
// The returned cleanup function should be called during application shutdown.
func InitTelemetry(ctx context.Context, mode Mode) (context.Context, func(context.Context), error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx = setLoggerInContext(ctx, logger)
	RecordInfoEvent(ctx, "Initializing telemetry")

	var cleanupFuncs []func(context.Context) error
	cleanup := func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second) // Reduced timeout for faster cleanup
		defer cancel()

		var wg sync.WaitGroup
		for idx := range cleanupFuncs {
			fn := cleanupFuncs[idx]
			wg.Go(func() {
				_ = fn(ctx)
			})
		}
		wg.Wait()
	}

	// Create OTEL resource once
	res, err := newResource(ctx)
	if err != nil {
		return nil, nil, err
	}
	RecordInfoEvent(ctx, "OTEL Resource created")

	// Create service configuration from resource attributes
	serviceConfig := newServiceConfig(res)
	if !serviceConfig.IsValid() {
		return nil, nil, fmt.Errorf("invalid service configuration: missing required attributes")
	}
	ctx = setServiceConfigInContext(ctx, serviceConfig)
	RecordInfoEvent(ctx, fmt.Sprintf("Service configuration initialized: %+v", serviceConfig))

	// Initialize log provider
	lp, err := newOTELLogProvider(ctx, res)
	if err != nil {
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, lp.Shutdown)
	logger = slog.New(otelslog.NewHandler(instrumentationIdentifier, otelslog.WithLoggerProvider(lp)))
	ctx = setLoggerInContext(ctx, logger)
	RecordInfoEvent(ctx, "Logger initialized")

	// Initialize trace provider
	tp, err := newOTELTraceProvider(ctx, res, mode)
	if err != nil {
		cleanup(ctx)
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, tp.Shutdown)
	otel.SetTracerProvider(otelpyroscope.NewTracerProvider(tp)) // Wrap TracerProvider so that profiling samples get span_id/trace_id labels
	RecordInfoEvent(ctx, "Tracer initialized")

	// Initialize metrics provider
	mp, err := newOTELMetricsProvider(ctx, res)
	if err != nil {
		cleanup(ctx)
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, mp.Shutdown)
	otel.SetMeterProvider(mp)
	RecordInfoEvent(ctx, "Meter initialized")

	// Initialize metrics collector after meter provider is set
	collector, err := NewMetricsCollector()
	if err != nil {
		cleanup(ctx)
		return nil, nil, err
	}
	ctx = setMetricsCollectorInContext(ctx, collector)
	RecordInfoEvent(ctx, "Metrics collector initialized")

	profiler, err := startProfiler(ctx, serviceConfig)
	if err != nil {
		RecordErrorEvent(ctx, err)
	} else {
		cleanupFuncs = append(cleanupFuncs, func(ctx context.Context) error {
			return profiler.Stop()
		})
	}

	RecordInfoEvent(ctx, "Telemetry initialization complete")

	return ctx, cleanup, nil
}

// Mode represents the telemetry/logging mode.
type Mode int

const (
	// ModeDebug enables debug monitoring and logging.
	ModeDebug Mode = iota
	// ModeProd enables production monitoring.
	ModeProd
)

func (m Mode) String() string {
	switch m {
	case ModeDebug:
		return "ModeDebug"
	case ModeProd:
		return "ModeProd"
	default:
		return "Unknown"
	}
}
