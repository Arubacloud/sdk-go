//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/securego/gosec/v2/cmd/gosec"
	_ "github.com/vektra/mockery/v2"
	_ "mvdan.cc/gofumpt"
)
