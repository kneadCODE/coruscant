package telemetry

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneCopiesLogger(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	ctx := setLoggerInContext(context.Background(), logger)
	cloned := Clone(ctx)
	assert.Equal(t, LoggerFromContext(ctx), LoggerFromContext(cloned))
}

func TestCloneEmptyContext(t *testing.T) {
	ctx := context.Background()
	cloned := Clone(ctx)
	assert.Nil(t, LoggerFromContext(cloned))
}
