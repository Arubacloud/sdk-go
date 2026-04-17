package restclient

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Arubacloud/sdk-go/internal/impl/logger/noop"
)

func newTestClient() *Client {
	return &Client{logger: &noop.NoOpLogger{}}
}

func TestWaitForResourceState(t *testing.T) {
	const fast = 1 * time.Millisecond

	t.Run("success on first attempt without initial sleep", func(t *testing.T) {
		calls := 0
		getter := func(_ context.Context) (string, error) {
			calls++
			return "Active", nil
		}
		cfg := PollingConfig{MaxAttempts: 5, Interval: 100 * time.Millisecond, SuccessStates: []string{"Active"}}

		start := time.Now()
		err := newTestClient().WaitForResourceState(context.Background(), "R", "id", getter, cfg)
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 getter call, got %d", calls)
		}
		// Should return well before one Interval (100ms). 50ms is a comfortable ceiling on any CI machine.
		if elapsed >= 50*time.Millisecond {
			t.Fatalf("first attempt took %v; initial sleep not removed", elapsed)
		}
	})

	t.Run("success after multiple attempts", func(t *testing.T) {
		calls := 0
		getter := func(_ context.Context) (string, error) {
			calls++
			if calls < 3 {
				return "Pending", nil
			}
			return "Active", nil
		}
		cfg := PollingConfig{MaxAttempts: 5, Interval: fast, SuccessStates: []string{"Active"}}

		if err := newTestClient().WaitForResourceState(context.Background(), "R", "id", getter, cfg); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if calls != 3 {
			t.Fatalf("expected 3 getter calls, got %d", calls)
		}
	})

	t.Run("failure state returns error immediately", func(t *testing.T) {
		getter := func(_ context.Context) (string, error) { return "Failed", nil }
		cfg := PollingConfig{MaxAttempts: 5, Interval: fast, SuccessStates: []string{"Active"}, FailureStates: []string{"Failed"}}

		err := newTestClient().WaitForResourceState(context.Background(), "R", "id", getter, cfg)
		if err == nil {
			t.Fatal("expected error for failure state")
		}
		if !strings.Contains(err.Error(), "reached failure state") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})

	t.Run("timeout preserves last known state", func(t *testing.T) {
		getter := func(_ context.Context) (string, error) { return "Pending", nil }
		cfg := PollingConfig{MaxAttempts: 3, Interval: fast, SuccessStates: []string{"Active"}}

		err := newTestClient().WaitForResourceState(context.Background(), "R", "id", getter, cfg)
		if err == nil {
			t.Fatal("expected timeout error")
		}
		if !strings.Contains(err.Error(), "last state: Pending") {
			t.Fatalf("expected last state in error, got: %v", err)
		}
	})

	t.Run("timeout wraps getter error when no state was ever observed", func(t *testing.T) {
		sentinel := errors.New("api unreachable")
		getter := func(_ context.Context) (string, error) { return "", sentinel }
		cfg := PollingConfig{MaxAttempts: 3, Interval: fast, SuccessStates: []string{"Active"}}

		err := newTestClient().WaitForResourceState(context.Background(), "R", "id", getter, cfg)
		if err == nil {
			t.Fatal("expected timeout error")
		}
		if !errors.Is(err, sentinel) {
			t.Fatalf("expected sentinel in error chain, got: %v", err)
		}
	})

	t.Run("context cancellation is respected before first attempt", func(t *testing.T) {
		calls := 0
		getter := func(_ context.Context) (string, error) {
			calls++
			return "Pending", nil
		}
		cfg := PollingConfig{MaxAttempts: 5, Interval: fast, SuccessStates: []string{"Active"}}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := newTestClient().WaitForResourceState(ctx, "R", "id", getter, cfg)
		if err == nil {
			t.Fatal("expected error on cancelled context")
		}
		if !strings.Contains(err.Error(), "context cancelled") {
			t.Fatalf("unexpected error message: %v", err)
		}
		if calls != 0 {
			t.Fatalf("getter should not have been called; called %d times", calls)
		}
	})
}
