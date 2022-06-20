package token

import "testing"

func TestTokenString(t *testing.T) {
	// string
	s := &Token{Type: STRING, Literal: "Moi"}
	if s.String() != "Moi" {
		t.Fatalf("Unexpected string-version of token.STRING")
	}

	// misc
	m := &Token{Type: SEMICOLON, Literal: ";"}
	if m.String() != "token{Type:; Literal:;}" {
		t.Fatalf("Unexpected string-version of token.SEMICOLON")
	}

}
