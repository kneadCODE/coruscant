package telemetry

import (
	"context"
	"fmt"
	"os"

	"github.com/grafana/pyroscope-go"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func startProfiler(ctx context.Context, sc ServiceConfig) (*pyroscope.Profiler, error) {
	RecordInfoEvent(ctx, "Starting Pyroscope profiler")

	profiler, err := pyroscope.Start(pyroscope.Config{
		// ApplicationName: sc.System + "/" + sc.Name,
		ApplicationName: sc.Name,
		Tags: map[string]string{
			replaceAttrKeyDotWithUnderscore(semconv.ServiceNameKey):               sc.Name,
			replaceAttrKeyDotWithUnderscore(semconv.ServiceNamespaceKey):          sc.System,
			replaceAttrKeyDotWithUnderscore(semconv.ServiceVersionKey):            sc.Version,
			replaceAttrKeyDotWithUnderscore(semconv.DeploymentEnvironmentNameKey): sc.Environment,
			replaceAttrKeyDotWithUnderscore(semconv.HostNameKey):                  sc.HostName,
		},
		ServerAddress: os.Getenv("PYROSCOPE_SERVER_ADDRESS"),
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
		// Logger:     &slogLogrusAdapter{logger: LoggerFromContext(ctx)},
		// UploadRate: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start pyroscope profiler: %w", err)
	}

	RecordInfoEvent(ctx, "Pyroscope profiler started successfully")
	return profiler, nil
}

// type slogLogrusAdapter struct {
// 	logger *slog.Logger
// }

// func (a *slogLogrusAdapter) Infof(format string, args ...interface{}) {
// 	a.logger.Info(fmt.Sprintf(format, args...))
// }

// func (a *slogLogrusAdapter) Debugf(format string, args ...interface{}) {
// 	a.logger.Debug(fmt.Sprintf(format, args...))
// }

// func (a *slogLogrusAdapter) Errorf(format string, args ...interface{}) {
// 	a.logger.Error(fmt.Sprintf(format, args...))
// }
