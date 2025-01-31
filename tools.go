//go:build tools

package tools

// This file purely exists to make sure that certain applications are available for generating
// Because of the go:build tools flag at the top, this file never gets compiled and the errors don't matter
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint" // the linter
	_ "github.com/ogen-go/ogen/cmd/ogen"                    // openapi spec generator
)
