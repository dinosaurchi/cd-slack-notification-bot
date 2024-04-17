//go:build tools
// +build tools

// Build tag reference: https://pkg.go.dev/cmd/go#hdr-Build_constraints

package main

import (
	_ "golang.org/x/lint/golint"
	_ "gotest.tools/gotestsum"
)
