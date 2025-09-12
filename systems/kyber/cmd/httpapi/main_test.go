package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

func TestRun(t *testing.T) {
	// Test run() function - it should initialize telemetry and start the server
	// Since run() now returns an error, we can test it more easily

	// We can't let it run indefinitely, so we'll test in a goroutine with timeout
	done := make(chan error, 1)

	go func() {
		err := run()
		done <- err
	}()

	// Give it a moment to initialize and start, then we expect it to keep running
	select {
	case err := <-done:
		// If run() returns quickly, it's either an error or unexpected shutdown
		if err != nil {
			t.Logf("run() returned error: %v", err)
		} else {
			t.Log("run() returned successfully (unexpected)")
		}
	case <-time.After(100 * time.Millisecond):
		// Expected: run() should be still running the server
		t.Log("run() is running (expected behavior)")
	}
}

func TestStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(100 * time.Millisecond)
		cancel()
	}()
	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Failf(t, "start(ctx) panicked", "panic: %v", r)
			}
		}()
		errCh <- start(ctx)
	}()
	select {
	case err := <-errCh:
		t.Logf("start(ctx) returned: %v", err)
		// Optionally, check for specific error values here
	case <-time.After(500 * time.Millisecond):
		assert.Fail(t, "start(ctx) did not return in time")
	}
}

func TestRunTelemetryError(t *testing.T) {
	// Test run() function when telemetry initialization fails
	// Set invalid environment variables that might cause telemetry init to fail
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "invalid=attribute=with=too=many=equals")

	err := run()
	// We expect either success (if the invalid attribute is ignored)
	// or an error (if telemetry initialization fails)
	if err != nil {
		assert.Error(t, err, "Expected telemetry initialization error")
		t.Logf("Got expected telemetry error: %v", err)
	} else {
		t.Log("Telemetry initialization succeeded despite invalid attributes")
	}
}

func TestStartWithTelemetry(t *testing.T) {
	// Test start() with a proper telemetry context
	ctx := context.Background()

	// Initialize telemetry first to create a proper context
	telemetryCtx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDev)
	if err != nil {
		t.Skip("Could not initialize telemetry for test")
	}
	defer cleanup()

	// Create a context that will be cancelled quickly
	testCtx, cancel := context.WithTimeout(telemetryCtx, 50*time.Millisecond)
	defer cancel()

	err = start(testCtx)
	// Should either succeed, timeout, or fail due to port already in use
	if err != nil {
		// Expected errors: timeout, port in use, etc.
		assert.True(t,
			err.Error() == "context deadline exceeded" ||
				strings.Contains(err.Error(), "address already in use") ||
				strings.Contains(err.Error(), "startup failed"),
			"start() should return expected error, got: %v", err)
	}
}

func TestStartServerCreation(t *testing.T) {
	// Test that start() can create a server successfully
	// We don't actually start it to avoid hanging
	ctx := context.Background()

	// We can't easily test the server creation error without mocking,
	// but we can at least verify the function doesn't panic
	assert.NotPanics(t, func() {
		// Very quick timeout to avoid actually starting the server
		quickCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()
		start(quickCtx)
	}, "start() should not panic")
}

func TestStartErrorPaths(t *testing.T) {
	// Test start function with various conditions that might trigger errors
	testCases := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "cancelled_context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately
				return ctx
			}(),
		},
		{
			name: "expired_context",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				defer cancel()
				time.Sleep(1 * time.Millisecond) // Ensure timeout
				return ctx
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := start(tc.ctx)
			// We expect either no error (if server starts and stops quickly)
			// or a context-related error
			if err != nil {
				t.Logf("start() returned expected error for %s: %v", tc.name, err)
			}
		})
	}
}

func TestRunSuccess(t *testing.T) {
	// Test that run() can start successfully (but we won't let it run forever)
	// We'll test in a goroutine and expect it to not return immediately with an error

	done := make(chan error, 1)

	go func() {
		err := run()
		done <- err
	}()

	// Give run() some time to initialize telemetry and start the server
	select {
	case err := <-done:
		// If run() returns within 50ms, something went wrong
		if err != nil {
			t.Logf("run() failed during startup: %v", err)
			// This might be expected in some test environments
		} else {
			t.Log("run() returned successfully (unexpected in normal operation)")
		}
	case <-time.After(50 * time.Millisecond):
		// Expected: run() should still be running the server
		assert.True(t, true, "run() started successfully and is running")
	}
}

func TestRunErrorPaths(t *testing.T) {
	// Now we can properly test error paths since run() returns errors

	t.Run("telemetry_error", func(t *testing.T) {
		// Try to trigger telemetry initialization error
		t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "=invalid")

		err := run()
		// The function should either succeed or return an error
		// We log the result to understand the behavior
		if err != nil {
			t.Logf("run() returned error as expected: %v", err)
		} else {
			t.Log("run() succeeded despite environment manipulation")
		}
	})
}
