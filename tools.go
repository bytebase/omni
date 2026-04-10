//go:build tools

// Package tools pins tool and upcoming-feature dependencies so that
// go mod tidy does not remove them before they are directly imported.
package tools

import _ "github.com/hjson/hjson-go/v4"
