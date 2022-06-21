// Package stdlib contains a simple/small standard-library, which
// is written in TCL itself.
//
package stdlib

import (
	_ "embed" // embedded-resource magic
)

//go:embed stdlib.tcl
var message string

// Contents returns the embedded TCL code.
func Contents() []byte {
	return []byte(message)
}
