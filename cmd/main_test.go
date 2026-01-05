package main

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

func TestRun_Logging(t *testing.T) {
	t.Run("Logging Enabled by Default", func(t *testing.T) {
		// Create a context that cancels quickly to stop the server
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		stdout := &bytes.Buffer{}
		stderr := io.Discard
		stdin := strings.NewReader("")
		args := []string{"cmd"}

		// Mock getenv to return defaults
		getenv := func(key, defaultValue string) string {
			return defaultValue
		}

		err := run(ctx, args, getenv, stdin, stdout, stderr)
		if err != nil {
			t.Fatalf("run failed: %v", err)
		}

		if !strings.Contains(stdout.String(), "listening") {
			t.Errorf("expected logs to contain 'listening', got: %q", stdout.String())
		}
	})

	t.Run("Logging Disabled via Env", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		stdout := &bytes.Buffer{}
		stderr := io.Discard
		stdin := strings.NewReader("")
		args := []string{"cmd"}

		// Mock getenv to simulate DISABLE_LOGGING=true
		getenv := func(key, defaultValue string) string {
			if key == "DISABLE_LOGGING" {
				return "true"
			}
			return defaultValue
		}

		err := run(ctx, args, getenv, stdin, stdout, stderr)
		if err != nil {
			t.Fatalf("run failed: %v", err)
		}

		if stdout.Len() > 0 {
			t.Errorf("expected no logs, got: %q", stdout.String())
		}
	})
}
