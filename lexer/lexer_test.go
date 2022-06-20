package lexer

import (
	"os"
	"strings"
	"testing"

	"github.com/skx/critical/token"
)

// TestEmpty tests a couple of different empty strings
func TestEmpty(t *testing.T) {

	// empty string
	l := New("")

	// should return EOF for N times
	i := 0

	for i < 5 {
		tok := l.NextToken()
		if tok.Type != token.EOF {
			t.Fatalf("expected EOF, got %v", tok)
		}

		p := l.peekChar()
		if p != rune(0) {
			t.Fatalf("peeking past EOF failed")
		}

		i++
	}
}

// TestVariable does simple variable testing.
func TestVariable(t *testing.T) {
	input := `$a+$b`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.VARIABLE, "$a"},
		{token.IDENT, "+"},
		{token.VARIABLE, "$b"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q: %v", i, tt.expectedType, tok.Type, tok)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q: %v", i, tt.expectedLiteral, tok.Literal, tok)
		}
	}
}

// TestEscape ensures that strings have escape-characters processed.
func TestStringEscape(t *testing.T) {
	input := `"Steve\n\r\\" "Kemp\n\t\n" "Inline \"quotes\"."`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.STRING, "Steve\n\r\\"},
		{token.STRING, "Kemp\n\t\n"},
		{token.STRING, "Inline \"quotes\"."},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q: %v", i, tt.expectedType, tok.Type, tok)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q: %v", i, tt.expectedLiteral, tok.Literal, tok)
		}
	}
}

// TestUnterminatedString ensures that an unclosed-string is an error
func TestUnterminatedString(t *testing.T) {
	input := `"Steve`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ILLEGAL, "unterminated string"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestContinue checks we continue newlines.
func TestContinue(t *testing.T) {
	input := `"This is a test \
which continues"`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.STRING, "This is a test which continues"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestParseNumber ensures we can catch errors in numbers
//
// This includes handling the unary +/- prefixes which might be present.
func TestParseNumber(t *testing.T) {

	// Parsing a number
	lex := New("449691189")
	tok := lex.NextToken()

	if tok.Type != token.NUMBER {
		t.Fatalf("parsed number as wrong type")
	}
	if tok.Literal != "449691189" {
		t.Fatalf("error lexing got:%s", tok.Literal)
	}

	// Now a malformed number
	lex = New("10-10")
	tok = lex.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("parsed number as wrong type")
	}
	if !strings.Contains(tok.Literal, "'-' may only occur at the start of the number") {
		t.Fatalf("got error, but wrong one: %s", tok.Literal)
	}

	// Now a number that's out of range.
	lex = New("18446744073709551620")

	tok = lex.NextToken()

	if tok.Type != token.ILLEGAL {
		t.Fatalf("parsed number as wrong type")
	}
	if !strings.Contains(tok.Literal, "out of range") {
		t.Fatalf("got error, but wrong one: %s", tok.Literal)
	}

}

// TestIllegal looks for some illegal inputs
func TestIllegal(t *testing.T) {
	illegals := []string{
		"{",
		"}",
		"[",
		"]",
		"\"steve",
	}

	for _, test := range illegals {
		lex := New(test)

		tok := lex.NextToken()

		if tok.Type != token.ILLEGAL {
			t.Fatalf("expected ILLEGAL, got %v", tok)
		}
	}
}

func TestNested(t *testing.T) {
	input := `[expr 1 + [expr 2+3]]`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.EVAL, "[expr 1 + [expr 2+3]]"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q: %v", i, tt.expectedType, tok.Type, tok)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q: %v", i, tt.expectedLiteral, tok.Literal, tok)
		}

	}
}

// TestInteger tests that we parse integers appropriately.
//
// This includes handling the unary +/- prefixes which might be present.
func TestInteger(t *testing.T) {

	old := os.Getenv("DECIMAL_NUMBERS")
	os.Setenv("DECIMAL_NUMBERS", "true")

	type TestCase struct {
		input  string
		output string
	}

	tests := []TestCase{
		{input: "3", output: "3"},
		{input: "-3", output: "-3"},
		{input: "-0", output: "0"},
		{input: "-10", output: "-10"},
		{input: "0xff", output: "255"},
		{input: "0b11111111", output: "255"},
	}

	for _, tst := range tests {

		lex := New(tst.input)
		tok := lex.NextToken()

		if tok.Type != token.NUMBER {
			t.Fatalf("failed to parse '%s' as number: %s", tst.input, tok)
		}
		if tok.Literal != tst.output {
			t.Fatalf("error lexing %s - expected:%s got:%s", tst.input, tst.output, tok.Literal)
		}
	}

	os.Setenv("DECIMAL_NUMBERS", old)
}

func TestEval(t *testing.T) {

	input := `puts [expr 1 + [ expr 2 + 3] ]`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "puts"},
		{token.EVAL, "[expr 1 + [ expr 2 + 3] ]"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q: %v", i, tt.expectedType, tok.Type, tok)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q: %v", i, tt.expectedLiteral, tok.Literal, tok)
		}
	}
}

func TestIf(t *testing.T) {

	input := `if { $i < 10 } { puts "OK" } else { puts "FAILED"}`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "if"},
		{token.BLOCK, " $i < 10 "},
		{token.BLOCK, " puts \"OK\" "},
		{token.IDENT, "else"},
		{token.BLOCK, " puts \"FAILED\""},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q: %v", i, tt.expectedType, tok.Type, tok)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q: %v", i, tt.expectedLiteral, tok.Literal, tok)
		}
	}
}
