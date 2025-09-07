package main

import (
	"context"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(100 * time.Millisecond)
		cancel()
	}()
	go func() {
		run(ctx)
		close(done)
	}()
	select {
	case <-done:
		// Success: run() returned
	case <-time.After(500 * time.Millisecond):
		t.Error("run() did not return in time")
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
				t.Errorf("start(ctx) panicked: %v", r)
			}
		}()
		errCh <- start(ctx)
	}()
	select {
	case err := <-errCh:
		t.Logf("start(ctx) returned: %v", err)
		// Optionally, check for specific error values here
	case <-time.After(500 * time.Millisecond):
		t.Error("start(ctx) did not return in time")
	}
}
