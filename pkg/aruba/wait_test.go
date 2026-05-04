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
// setRefresh / setTerminalStates setters
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

func TestStatusMixin_SetTerminalStates(t *testing.T) {
	var m statusMixin
	ts := map[string]bool{"Active": true, "Error": false}
	m.setTerminalStates(ts)
	if len(m.terminalStates) != 2 {
		t.Errorf("terminalStates len = %d, want 2", len(m.terminalStates))
	}
	if !m.terminalStates["Active"] {
		t.Error("terminalStates[Active] should be true")
	}
	if m.terminalStates["Error"] {
		t.Error("terminalStates[Error] should be false")
	}
}

// --------------------------------------------------------------------------
// WaitUntilState / WaitUntilActive — nil refresh
// --------------------------------------------------------------------------

func TestWaitUntilState_RefreshNil_Error(t *testing.T) {
	var m statusMixin
	err := m.WaitUntilState(context.Background(), "Active")
	if err == nil {
		t.Fatal("expected error when refresh is nil")
	}
	if !strings.Contains(err.Error(), "refresh callback not set") {
		t.Errorf("error message = %q; want 'refresh callback not set'", err.Error())
	}
}

// --------------------------------------------------------------------------
// WaitUntilActive / WaitUntilState — happy paths
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
	state := "Pending"
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls >= 3 {
			state = "Active"
		}
		s := state
		m.setStatus(&types.ResourceStatus{State: &s})
		return nil
	})
	if err := m.WaitUntilActive(context.Background(), fastOpts()...); err != nil {
		t.Fatalf("WaitUntilActive error: %v", err)
	}
	if m.State() != "Active" {
		t.Errorf("State() = %q after wait, want Active", m.State())
	}
	if calls < 3 {
		t.Errorf("refresh called %d times, want >= 3", calls)
	}
}

func TestWaitUntilState_CustomTarget(t *testing.T) {
	var m statusMixin
	calls := 0
	state := "Pending"
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls >= 2 {
			state = "Available"
		}
		s := state
		m.setStatus(&types.ResourceStatus{State: &s})
		return nil
	})
	if err := m.WaitUntilState(context.Background(), "Available", fastOpts()...); err != nil {
		t.Fatalf("WaitUntilState error: %v", err)
	}
	if m.State() != "Available" {
		t.Errorf("State() = %q after wait, want Available", m.State())
	}
}

// --------------------------------------------------------------------------
// Error terminal state
// --------------------------------------------------------------------------

func TestWaitUntilActive_ErrorTerminal(t *testing.T) {
	var m statusMixin
	m.setTerminalStates(map[string]bool{"Active": true, "Error": false})
	state := "Pending"
	calls := 0
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls >= 2 {
			state = "Error"
		}
		s := state
		m.setStatus(&types.ResourceStatus{State: &s})
		return nil
	})
	err := m.WaitUntilActive(context.Background(), fastOpts()...)
	if err == nil {
		t.Fatal("expected error when terminal error state reached")
	}
	if !strings.Contains(err.Error(), "terminal error state") {
		t.Errorf("error = %q; want 'terminal error state'", err.Error())
	}
}

// --------------------------------------------------------------------------
// Transient refresh errors are retried
// --------------------------------------------------------------------------

func TestWaitUntilActive_RefreshError_Retried(t *testing.T) {
	var m statusMixin
	calls := 0
	state := "Pending"
	m.setRefresh(func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("transient network error")
		}
		state = "Active"
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
	state := "Pending"
	s := state
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
	state := "Pending"
	s := state
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
	state := "Pending"
	s := state
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
