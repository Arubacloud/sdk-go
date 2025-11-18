package restclient

import (
	"fmt"
	"net/http"
)

// RequestEditorFn is a function that can modify an HTTP request before it's sent
// This is compatible with oapi-codegen's RequestEditorFn
type RequestEditorFn func(req *http.Request) error

// WithCustomHeaders returns a RequestEditorFn that adds custom headers
func WithCustomHeaders(headers map[string]string) RequestEditorFn {
	return func(req *http.Request) error {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return nil
	}
}

// WithUserAgent returns a RequestEditorFn that sets the User-Agent header
func WithUserAgent(userAgent string) RequestEditorFn {
	return func(req *http.Request) error {
		req.Header.Set("User-Agent", userAgent)
		return nil
	}
}

// WithContentType returns a RequestEditorFn that sets the Content-Type header
func WithContentType(contentType string) RequestEditorFn {
	return func(req *http.Request) error {
		req.Header.Set("Content-Type", contentType)
		return nil
	}
}

// WithAccept returns a RequestEditorFn that sets the Accept header
func WithAccept(accept string) RequestEditorFn {
	return func(req *http.Request) error {
		req.Header.Set("Accept", accept)
		return nil
	}
}

// ChainRequestEditors chains multiple RequestEditorFn together
func ChainRequestEditors(editors ...RequestEditorFn) RequestEditorFn {
	return func(req *http.Request) error {
		for _, editor := range editors {
			if err := editor(req); err != nil {
				return fmt.Errorf("request editor failed: %w", err)
			}
		}
		return nil
	}
}
