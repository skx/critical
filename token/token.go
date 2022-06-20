// Package token contains the token-types which our lexer produces,
// and which our parser understands.
package token

import "fmt"

// Type is a string.
type Type string

// Token struct represent the token which is returned from the lexer.
type Token struct {
	Type    Type
	Literal string
}

// pre-defined TokenTypes
const (
	// Things
	SEMICOLON = ";"
	EOF       = "EOF"
	NEWLINE   = "\n"

	// types
	BLOCK    = "BLOCK"
	EVAL     = "EVAL"
	IDENT    = "IDENT"
	ILLEGAL  = "ILLEGAL"
	NUMBER   = "NUMBER"
	STRING   = "STRING"
	VARIABLE = "VARIABLE"
)

// String turns the token into a readable string
func (t Token) String() string {

	// string?
	if t.Type == STRING {
		return t.Literal
	}

	// everything else is less pretty
	return fmt.Sprintf("token{Type:%s Literal:%s}", t.Type, t.Literal)
}
