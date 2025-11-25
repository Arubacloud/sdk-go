// Package interceptor provides basic interfaces and types for implementing
// request interception logic, typically used before send or processing an
// HTTP request.
package interceptor

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrInvalidInterceptFunc = errors.New("invalid intercept function")
	ErrInvalidHTTPRequest   = errors.New("invalid http request")
	ErrInterceptFuncFailed  = errors.New("intercept function failed")
)

// InterceptFunc is a function signature that defines the core interception
// logic.
//
// It receives a context.Context and the *http.Request.
//
// Implementations should return an error if the interception logic fails or if
// the request should be halted.
type InterceptFunc func(ctx context.Context, r *http.Request) error

// Interceptable is an interface implemented by components that want to have
// one or more InterceptFuncs bound to them, usually for execution before a
// core operation.
type Interceptable interface {
	// Bind adds the provided InterceptFuncs to the interceptable component's
	// execution chain.
	Bind(interceptFuncs ...InterceptFunc) error
}

// Interceptor is the core interface for executing the request interception
// logic.
//
// Components that implement this interface are responsible for running the
// bound InterceptFuncs.
type Interceptor interface {
	// Intercept executes all bound InterceptFuncs in order.
	//
	// It returns an error immediately upon the first InterceptFunc that fails.
	Intercept(ctx context.Context, request *http.Request) error
}
