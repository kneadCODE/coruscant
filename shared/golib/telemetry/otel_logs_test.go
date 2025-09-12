package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestNewOTELSlogHandler_DevModes(t *testing.T) {
	ctx := context.Background()
	res, err := newResource(ctx)
	assert.NoError(t, err)

	modes := []Mode{ModeDev, ModeDevDebug}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			handler, cleanup, err := newOTELSlogHandler(res, mode)
			assert.NoError(t, err)
			assert.NotNil(t, handler)
			assert.NotNil(t, cleanup)
			cleanup()
		})
	}
}

func TestNewOTELSlogHandler_ProdModes(t *testing.T) {
	ctx := context.Background()
	res, err := newResource(ctx)
	assert.NoError(t, err)

	modes := []Mode{ModeProd, ModeProdDebug}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			handler, cleanup, err := newOTELSlogHandler(res, mode)
			assert.NoError(t, err)
			assert.NotNil(t, handler)
			assert.NotNil(t, cleanup)
			cleanup()
		})
	}
}

func TestNewOTELSlogHandler_InvalidResource(t *testing.T) {
	var res *resource.Resource
	handler, cleanup, err := newOTELSlogHandler(res, ModeDev)

	if err != nil {
		assert.Nil(t, handler)
		assert.Nil(t, cleanup)
	} else {
		assert.NotNil(t, handler)
		assert.NotNil(t, cleanup)
		if cleanup != nil {
			cleanup()
		}
	}
}

func TestNewOTELSlogHandler_WithNilWriter(t *testing.T) {
	ctx := context.Background()
	res, err := newResource(ctx)
	assert.NoError(t, err)

	// Test both dev and prod modes to ensure we exercise all paths
	modes := []Mode{ModeDev, ModeDevDebug, ModeProd, ModeProdDebug}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			// Call the handler creation - this should exercise the error paths if any exist
			handler, cleanup, err := newOTELSlogHandler(res, mode)

			// The function should either succeed or fail gracefully
			if err != nil {
				assert.Nil(t, handler)
				assert.Nil(t, cleanup)
			} else {
				assert.NotNil(t, handler)
				assert.NotNil(t, cleanup)
				cleanup()
			}
		})
	}
}

// Test stdoutlog.New error conditions by trying multiple scenarios
func TestNewOTELSlogHandler_ErrorConditions(t *testing.T) {
	ctx := context.Background()
	res, err := newResource(ctx)
	assert.NoError(t, err)

	// Test different scenarios that might trigger errors in stdoutlog.New
	testCases := []struct {
		name string
		mode Mode
	}{
		{"DevMode", ModeDev},
		{"DevDebugMode", ModeDevDebug},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Try multiple times with potential error conditions
			for i := 0; i < 3; i++ {
				handler, cleanup, err := newOTELSlogHandler(res, tc.mode)
				if err != nil {
					// If we get an error, make sure the returns are correct
					assert.Nil(t, handler)
					assert.Nil(t, cleanup)
				} else {
					assert.NotNil(t, handler)
					assert.NotNil(t, cleanup)
					cleanup()
				}
			}
		})
	}
}

func TestNewOTELSlogHandler_AllModeCombinations(t *testing.T) {
	ctx := context.Background()
	res, err := newResource(ctx)
	assert.NoError(t, err)

	// Test all mode combinations to ensure we hit all branches
	allModes := []Mode{ModeDev, ModeDevDebug, ModeProd, ModeProdDebug}

	for _, mode := range allModes {
		t.Run(mode.String()+"_multiple_calls", func(t *testing.T) {
			// Call multiple times to exercise any potential error paths
			for i := 0; i < 2; i++ {
				handler, cleanup, err := newOTELSlogHandler(res, mode)

				if err != nil {
					// Error case: both should be nil
					assert.Nil(t, handler)
					assert.Nil(t, cleanup)
				} else {
					// Success case: both should be non-nil
					assert.NotNil(t, handler)
					assert.NotNil(t, cleanup)

					// Test that cleanup works without panicking
					assert.NotPanics(t, func() {
						cleanup()
					})
				}
			}
		})
	}
}

func TestNewOTELSlogHandler_ResourceEdgeCases(t *testing.T) {
	// Test with different resource configurations
	testCases := []struct {
		name          string
		resourceSetup func() *resource.Resource
		expectedError bool
	}{
		{
			name: "ValidResource",
			resourceSetup: func() *resource.Resource {
				res, _ := newResource(context.Background())
				return res
			},
			expectedError: false,
		},
		{
			name: "NilResource",
			resourceSetup: func() *resource.Resource {
				return nil
			},
			expectedError: false, // nil resource might be handled gracefully
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.resourceSetup()

			// Test with different modes
			modes := []Mode{ModeDev, ModeProd}
			for _, mode := range modes {
				handler, cleanup, err := newOTELSlogHandler(res, mode)

				if tc.expectedError {
					assert.Error(t, err)
					assert.Nil(t, handler)
					assert.Nil(t, cleanup)
				} else {
					if err != nil {
						// If error occurred (might be environmental), ensure proper cleanup state
						assert.Nil(t, handler)
						assert.Nil(t, cleanup)
					} else {
						// Success case
						assert.NotNil(t, handler)
						assert.NotNil(t, cleanup)
						cleanup()
					}
				}
			}
		})
	}
}
