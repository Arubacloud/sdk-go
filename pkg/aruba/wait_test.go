package aruba

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/async"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// refreshIsSet is a package-internal helper used by wrapper tests to assert
// that the adapter has injected a refresh callback.
func refreshIsSet(m *statusMixin) bool { return m.refresh != nil }

// --------------------------------------------------------------------------
// WaitOption machinery
// --------------------------------------------------------------------------

func TestWaitOption_Defaults(t *testing.T) {
	o := applyWaitOptions(nil)
	if o.retries != async.DefaultRetries {
		t.Errorf("retries = %d, want %d", o.retries, async.DefaultRetries)
	}
	if o.baseDelay != async.DefaultBaseDelay {
		t.Errorf("baseDelay = %v, want %v", o.baseDelay, async.DefaultBaseDelay)
	}
	if o.timeout != async.DefaultTimeout {
		t.Errorf("timeout = %v, want %v", o.timeout, async.DefaultTimeout)
	}
}

func TestWaitOption_WithRetries(t *testing.T) {
	o := applyWaitOptions([]WaitOption{WithRetries(5)})
	if o.retries != 5 {
		t.Errorf("retries = %d, want 5", o.retries)
	}
}

func TestWaitOption_WithBaseDelay(t *testing.T) {
	d := 42 * time.Millisecond
	o := applyWaitOptions([]WaitOption{WithBaseDelay(d)})
	if o.baseDelay != d {
		t.Errorf("baseDelay = %v, want %v", o.baseDelay, d)
	}
}

func TestWaitOption_WithTimeout(t *testing.T) {
	d := 7 * time.Second
	o := applyWaitOptions([]WaitOption{WithTimeout(d)})
	if o.timeout != d {
		t.Errorf("timeout = %v, want %v", o.timeout, d)
	}
}

func TestWaitOption_NilOption_Skipped(t *testing.T) {
	o := applyWaitOptions([]WaitOption{nil, WithRetries(3), nil})
	if o.retries != 3 {
		t.Errorf("retries = %d, want 3", o.retries)
	}
}

// --------------------------------------------------------------------------
// setRefresh setter
// --------------------------------------------------------------------------

func TestStatusMixin_SetRefresh(t *testing.T) {
	var m statusMixin
	called := false
	fn := func(_ context.Context) error { called = true; return nil }
	m.setRefresh(fn)
	if m.refresh == nil {
		t.Fatal("setRefresh: refresh field is still nil")
	}
	_ = m.refresh(context.Background())
	if !called {
		t.Error("setRefresh: injected function was not called")
	}
}

// --------------------------------------------------------------------------
// WaitUntilStates / WaitUntilActive — nil refresh
// --------------------------------------------------------------------------

func TestWaitUntilStates_RefreshNil_Error(t *testing.T) {
	var m statusMixin
	err := m.WaitUntilStates(context.Background(), []types.State{types.StateActive})
	if err == nil {
		t.Fatal("expected error when refresh is nil")
	}
	if !strings.Contains(err.Error(), "refresh callback not set") {
		t.Errorf("error message = %q; want 'refresh callback not set'", err.Error())
	}
}

// --------------------------------------------------------------------------
// WaitUntilActive / WaitUntilStates / WaitUntilReady — happy paths
// --------------------------------------------------------------------------

func fastOpts() []WaitOption {
	return []WaitOption{
		WithRetries(20),
		WithBaseDelay(1 * time.Millisecond),
		WithTimeout(2 * time.Second),
	}
}

func TestWaitUntilActive_HappyPath(t *testing.T) {
	var m statusMixin
	calls := 0
	var state types.State = "InCreation"
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls >= 3 {
			state = types.StateActive
		}
		s := state
		m.setStatus(&types.ResourceStatus{State: &s})
		return nil
	})
	if err := m.WaitUntilActive(context.Background(), fastOpts()...); err != nil {
		t.Fatalf("WaitUntilActive error: %v", err)
	}
	if m.State() != types.StateActive {
		t.Errorf("State() = %q after wait, want Active", m.State())
	}
	if calls < 3 {
		t.Errorf("refresh called %d times, want >= 3", calls)
	}
}

func TestWaitUntilStates_CustomTarget(t *testing.T) {
	var m statusMixin
	calls := 0
	var state types.State = "InCreation"
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls >= 2 {
			state = "Available"
		}
		s := state
		m.setStatus(&types.ResourceStatus{State: &s})
		return nil
	})
	if err := m.WaitUntilStates(context.Background(), []types.State{"Available"}, fastOpts()...); err != nil {
		t.Fatalf("WaitUntilStates error: %v", err)
	}
	if m.State() != "Available" {
		t.Errorf("State() = %q after wait, want Available", m.State())
	}
}

func TestWaitUntilReady_HappyPath(t *testing.T) {
	for _, target := range []types.State{
		types.StateActive, types.StateRunning, types.StateStopped,
		types.StateNotUsed, types.StateReserved, types.StateInUse, types.StateUsed,
	} {
		target := target
		t.Run(string(target), func(t *testing.T) {
			var m statusMixin
			calls := 0
			state := types.StateInCreation
			m.setRefresh(func(_ context.Context) error {
				calls++
				if calls >= 2 {
					state = target
				}
				s := state
				m.setStatus(&types.ResourceStatus{State: &s})
				return nil
			})
			if err := m.WaitUntilReady(context.Background(), fastOpts()...); err != nil {
				t.Fatalf("WaitUntilReady error for target %q: %v", target, err)
			}
			if m.State() != target {
				t.Errorf("State() = %q after wait, want %q", m.State(), target)
			}
		})
	}
}

// --------------------------------------------------------------------------
// Failure state — fast fail
// --------------------------------------------------------------------------

func TestWaitUntilActive_FailureState(t *testing.T) {
	for _, failState := range []types.State{types.StateFailed, types.StateError, types.StateDisabled} {
		failState := failState
		t.Run(string(failState), func(t *testing.T) {
			var m statusMixin
			calls := 0
			state := types.StateInCreation
			m.setRefresh(func(_ context.Context) error {
				calls++
				if calls >= 2 {
					state = failState
				}
				s := state
				m.setStatus(&types.ResourceStatus{State: &s})
				return nil
			})
			err := m.WaitUntilActive(context.Background(), fastOpts()...)
			if err == nil {
				t.Fatal("expected error when failure state reached")
			}
			if !strings.Contains(err.Error(), "failure state") {
				t.Errorf("error = %q; want 'failure state'", err.Error())
			}
		})
	}
}

// --------------------------------------------------------------------------
// Settled non-target state — fast fail (rule 4)
// --------------------------------------------------------------------------

func TestWaitUntilActive_SettledNonTarget(t *testing.T) {
	var m statusMixin
	state := types.StateInCreation
	calls := 0
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls >= 2 {
			state = types.StateReserved // settled non-target for WaitUntilActive
		}
		s := state
		m.setStatus(&types.ResourceStatus{State: &s})
		return nil
	})
	err := m.WaitUntilActive(context.Background(), fastOpts()...)
	if err == nil {
		t.Fatal("expected error when resource settles in non-target state")
	}
	if !strings.Contains(err.Error(), "settled in state") {
		t.Errorf("error = %q; want 'settled in state'", err.Error())
	}
}

// --------------------------------------------------------------------------
// Reserved is a valid WaitUntilUsed target
// --------------------------------------------------------------------------

func TestWaitUntilStates_ReservedIsTarget(t *testing.T) {
	var m statusMixin
	state := types.StateInCreation
	calls := 0
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls >= 2 {
			state = types.StateReserved
		}
		s := state
		m.setStatus(&types.ResourceStatus{State: &s})
		return nil
	})
	// Listing Reserved as a target (as WaitUntilUsed does) must succeed.
	err := m.WaitUntilStates(context.Background(),
		[]types.State{types.StateInUse, types.StateUsed, types.StateReserved},
		fastOpts()...)
	if err != nil {
		t.Fatalf("WaitUntilStates with Reserved as target: unexpected error: %v", err)
	}
	if m.State() != types.StateReserved {
		t.Errorf("State() = %q, want Reserved", m.State())
	}
}

// WaitUntilNotUsed must fail fast when resource settles in Reserved.
func TestWaitUntilStates_ReservedFailsFastWhenNotTarget(t *testing.T) {
	var m statusMixin
	state := types.StateInCreation
	calls := 0
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls >= 2 {
			state = types.StateReserved
		}
		s := state
		m.setStatus(&types.ResourceStatus{State: &s})
		return nil
	})
	// NotUsed only — Reserved is not in the target list.
	err := m.WaitUntilStates(context.Background(),
		[]types.State{types.StateNotUsed},
		fastOpts()...)
	if err == nil {
		t.Fatal("expected error: Reserved is settled but not in target list")
	}
	if !strings.Contains(err.Error(), "settled in state") {
		t.Errorf("error = %q; want 'settled in state'", err.Error())
	}
}

// --------------------------------------------------------------------------
// Transient refresh errors are retried
// --------------------------------------------------------------------------

func TestWaitUntilActive_RefreshError_Retried(t *testing.T) {
	var m statusMixin
	calls := 0
	state := types.StateInCreation
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("transient network error")
		}
		state = types.StateActive
		s := state
		m.setStatus(&types.ResourceStatus{State: &s})
		return nil
	})
	if err := m.WaitUntilActive(context.Background(), fastOpts()...); err != nil {
		t.Fatalf("WaitUntilActive error: %v (expected transient errors to be retried)", err)
	}
}

// --------------------------------------------------------------------------
// Retries exhausted
// --------------------------------------------------------------------------

func TestWaitUntilActive_RetriesExhausted(t *testing.T) {
	var m statusMixin
	s := types.StateInCreation
	m.setStatus(&types.ResourceStatus{State: &s})
	m.setRefresh(func(_ context.Context) error { return nil }) // never advances state
	err := m.WaitUntilActive(context.Background(),
		WithRetries(2),
		WithBaseDelay(1*time.Millisecond),
		WithTimeout(5*time.Second),
	)
	if err == nil {
		t.Fatal("expected error after retries exhausted")
	}
}

// --------------------------------------------------------------------------
// Timeout
// --------------------------------------------------------------------------

func TestWaitUntilActive_Timeout(t *testing.T) {
	var m statusMixin
	s := types.StateInCreation
	m.setStatus(&types.ResourceStatus{State: &s})
	m.setRefresh(func(_ context.Context) error { return nil }) // never advances
	err := m.WaitUntilActive(context.Background(),
		WithRetries(1000),
		WithBaseDelay(1*time.Millisecond),
		WithTimeout(50*time.Millisecond),
	)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

// --------------------------------------------------------------------------
// Context cancellation
// --------------------------------------------------------------------------

func TestWaitUntilActive_ContextCancellation(t *testing.T) {
	var m statusMixin
	s := types.StateInCreation
	m.setStatus(&types.ResourceStatus{State: &s})
	m.setRefresh(func(_ context.Context) error { return nil })
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err := m.WaitUntilActive(ctx, fastOpts()...)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
