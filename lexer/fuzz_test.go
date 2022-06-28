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
	f.Add([]byte(`set a 3.14`))
	f.Add([]byte(`set a 0xff`))
	f.Add([]byte(`set a 0b01010`))
	f.Add([]byte(`set a`))
	f.Add([]byte(`let b "Hello"`))

	// Errors
	f.Add([]byte(`set a "steve`))
	f.Add([]byte(`set a 10-21`))

	f.Fuzz(func(t *testing.T, input []byte) {

		// Create a new lexer
		c := New(string(input))

		var tok token.Token

		for tok.Type != token.EOF && tok.Type != token.ILLEGAL {
			tok = c.NextToken()
		}
	})
}
