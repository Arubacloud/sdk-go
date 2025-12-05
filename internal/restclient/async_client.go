package restclient

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	Response *http.Response
	Error    error
}

// Future represents a promise/future of an async call
type AsyncClient struct {
	resultCh chan Result
}

// Await blocks until the async operation finishes or context is canceled
func (f *AsyncClient) Await(ctx context.Context) (*http.Response, error) {
	select {
	case res := <-f.resultCh:
		return res.Response, res.Error
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// AsyncCallDefault triggers a function asynchronously and returns a Future
// with default retries and base delay values
func AsyncCallDefault(
	ctx context.Context,
	call func(ctx context.Context) (*http.Response, error),
) *AsyncClient {
	// default retries = 5, baseDelay = 500ms
	return AsyncCall(ctx, call, 5, 5*time.Second)
}

// AsyncCall triggers a function asynchronously and returns a Future
func AsyncCall(
	ctx context.Context,
	call func(ctx context.Context) (*http.Response, error),
	retries int,
	baseDelay time.Duration,
) *AsyncClient {
	fut := &AsyncClient{resultCh: make(chan Result, 1)}

	go func() {
		var lastErr error
		delay := baseDelay

		for attempt := 0; attempt <= retries; attempt++ {
			data, err := call(ctx)
			if err == nil {
				fut.resultCh <- Result{Response: data, Error: nil}
				return
			}

			lastErr = err

			select {
			case <-ctx.Done():
				fut.resultCh <- Result{Response: nil, Error: ctx.Err()}
				return
			case <-time.After(delay):
				// exponential backoff
				delay *= 2
			}
		}

		// failed after all retries
		fut.resultCh <- Result{Response: nil, Error: fmt.Errorf("after %d retries: %w", retries, lastErr)}
	}()

	return fut
}
