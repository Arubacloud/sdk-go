package restclient

import (
	"context"
	"fmt"
	"slices"
	"time"
)

// ResourceStateGetter is a function type that retrieves a resource and returns its state
type ResourceStateGetter func(ctx context.Context) (state string, err error)

// PollingConfig configures the resource state polling behavior
type PollingConfig struct {
	// MaxAttempts is the maximum number of polling attempts (default: 30)
	MaxAttempts int
	// Interval is the time between polling attempts (default: 5s)
	Interval time.Duration
	// SuccessStates are the states that indicate success (default: ["Active"])
	SuccessStates []string
	// FailureStates are the states that indicate failure (default: ["Failed", "Error"])
	FailureStates []string
}

// DefaultPollingConfig returns the default polling configuration
func DefaultPollingConfig() PollingConfig {
	return PollingConfig{
		MaxAttempts:   30,
		Interval:      5 * time.Second,
		SuccessStates: []string{"Active"},
		FailureStates: []string{"Failed", "Error"},
	}
}

// WaitForResourceState polls a resource until it reaches a success or failure state
func (c *Client) WaitForResourceState(ctx context.Context, resourceType, resourceID string, getter ResourceStateGetter, config PollingConfig) error {
	if config.MaxAttempts == 0 {
		config.MaxAttempts = 30
	}
	if config.Interval == 0 {
		config.Interval = 5 * time.Second
	}
	if len(config.SuccessStates) == 0 {
		config.SuccessStates = []string{"Active"}
	}
	if len(config.FailureStates) == 0 {
		config.FailureStates = []string{"Failed", "Error"}
	}

	c.logger.Debugf("Waiting for %s '%s' to become active...", resourceType, resourceID)

	var lastState string
	var lastErr error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for %s '%s'", resourceType, resourceID)
		default:
		}

		state, err := getter(ctx)
		if err != nil {
			lastErr = err
			c.logger.Debugf("Error checking %s '%s' status (attempt %d/%d): %v", resourceType, resourceID, attempt, config.MaxAttempts, err)
		} else {
			lastState = state
			lastErr = nil
			c.logger.Debugf("%s '%s' state: %s (attempt %d/%d)", resourceType, resourceID, state, attempt, config.MaxAttempts)

			if slices.Contains(config.SuccessStates, state) {
				c.logger.Debugf("%s '%s' is now %s", resourceType, resourceID, state)
				return nil
			}

			if slices.Contains(config.FailureStates, state) {
				return fmt.Errorf("%s '%s' reached failure state: %s", resourceType, resourceID, state)
			}
		}

		if attempt < config.MaxAttempts {
			time.Sleep(config.Interval)
		}
	}

	if lastState != "" {
		return fmt.Errorf("timeout waiting for %s '%s' to become active (last state: %s)", resourceType, resourceID, lastState)
	}
	if lastErr != nil {
		return fmt.Errorf("timeout waiting for %s '%s' to become active: %w", resourceType, resourceID, lastErr)
	}
	return fmt.Errorf("timeout waiting for %s '%s' to become active", resourceType, resourceID)
}
