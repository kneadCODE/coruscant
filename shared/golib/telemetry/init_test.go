package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTelemetry_AllModes(t *testing.T) {
	modes := []Mode{ModeDev, ModeDevDebug, ModeProd, ModeProdDebug}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			ctx := InitTelemetry(context.Background(), mode)
			logger := LoggerFromContext(ctx)
			assert.NotNil(t, logger)
		})
	}
}

func (m Mode) String() string {
	switch m {
	case ModeDev:
		return "ModeDev"
	case ModeDevDebug:
		return "ModeDevDebug"
	case ModeProd:
		return "ModeProd"
	case ModeProdDebug:
		return "ModeProdDebug"
	default:
		return "Unknown"
	}
}
