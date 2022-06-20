//go:build go1.18
// +build go1.18

package lexer

import (
	"testing"

	"github.com/skx/critical/token"
)

func FuzzLexer(f *testing.F) {

	// empty string
	f.Add([]byte(""))

	// simple entries
	f.Add([]byte("puts \"ok\""))
	f.Add([]byte("puts [expr 3 + 3]"))

	// Assignments
	f.Add([]byte(`set a 3`))
	f.Add([]byte(`set a`))
	f.Add([]byte(`let b "Hello"`))

	// Known errors are listed here.
	//
	// The purpose of fuzzing is to find panics, or unexpected errors.
	//
	// Some programs are obviously invalid though, so we don't want to
	// report those known-bad things.
	//	known := []string{}

	f.Fuzz(func(t *testing.T, input []byte) {

		// Create a new lexer
		c := New(string(input))

		var tok token.Token

		for tok.Type != token.EOF && tok.Type != token.ILLEGAL {
			tok = c.NextToken()
		}
	})
}
