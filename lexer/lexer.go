// Package lexer contains a simple lexer for reading an input-string
// and converting it into a series of tokens.
//
// In terms of syntax we're not very complex, so our lexer only needs
// to care about simple tokens:
//
// - Comments
// - Strings
// - Some simple characters such as "(", ")", "[", "]", "=>", "=", etc.
// -
//
// We can catch some basic errors in the lexing stage, such as unterminated
// strings, but the parser is the better place to catch such things.
package lexer

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/skx/critical/token"
)

// Lexer is used as the lexer for our deployr "language".
type Lexer struct {
	debug        bool                 // dump tokens as they're read?
	position     int                  // current character position
	readPosition int                  // next character position
	ch           rune                 // current character
	characters   []rune               // rune slice of input string
	lookup       map[rune]token.Token // lookup map for simple tokens
}

// New a Lexer instance from string input.
func New(input string) *Lexer {
	l := &Lexer{
		characters: []rune(input),
		debug:      false,
		lookup:     make(map[rune]token.Token),
	}
	l.readChar()

	if os.Getenv("DEBUG_LEXER") == "true" {
		l.debug = true
	}

	//
	// Lookup map of simple token-types.
	//
	l.lookup['\n'] = token.Token{Literal: "\\n", Type: token.NEWLINE}
	l.lookup[';'] = token.Token{Literal: ";", Type: token.SEMICOLON}
	l.lookup[rune(0)] = token.Token{Literal: "", Type: token.EOF}

	return l
}

// read one forward character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.characters) {
		l.ch = rune(0)
	} else {
		l.ch = l.characters[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) rewind() {
	if l.readPosition > 0 {
		l.position--
		l.readPosition--
		l.ch = l.characters[l.readPosition]
	}
}

// NextToken consumes and returns the next token from our input.
//
// It is a simple method which can optionally dump the tokens to the console
// if $DEBUG_LEXER is non-empty.
func (l *Lexer) NextToken() token.Token {

	tok := l.nextTokenReal()
	if l.debug {
		fmt.Printf("%v\n", tok)
	}

	return tok
}

// nextTokenReal does the real work of consuming and returning the next
// token from our input string.
func (l *Lexer) nextTokenReal() token.Token {
	var tok token.Token
	l.skipWhitespace()

	// Was this a simple token-type?
	val, ok := l.lookup[l.ch]
	if ok {
		// Yes, then skip the character itself, and return the
		// value we found.
		l.readChar()
		return val

	}

	// OK it wasn't a simple type
	switch l.ch {

	case rune(']'):
		tok.Type = token.ILLEGAL
		tok.Literal = "Closing ']' without opening one"

	case rune('}'):
		tok.Type = token.ILLEGAL
		tok.Literal = "Closing '}' without opening one"

	case rune('$'):
		val := l.readVariable()
		tok.Type = token.VARIABLE
		tok.Literal = val
	case rune('"'):
		str, err := l.readString()

		if err == nil {
			tok.Type = token.STRING
			tok.Literal = str
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = err.Error()
		}
	case rune('['):
		str, err := l.readEval()

		if err == nil {
			tok.Type = token.EVAL
			tok.Literal = "[" + str + "]"
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = err.Error()
		}
	case rune('{'):
		str, err := l.readBlock()

		if err == nil {
			tok.Type = token.BLOCK
			tok.Literal = str
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = err.Error()
		}
	default:
		// is it a number?
		if l.ch == '-' || isDigit(l.ch) {
			// Read it.
			tok = l.readDecimal()
			return tok
		}

		// is it an ident?
		tok.Literal = l.readIdentifier()
		tok.Type = token.IDENT

		return tok
	}

	// skip the character we've processed, and return the value
	l.readChar()
	return tok
}

// readDecimal returns a token consisting of decimal numbers, base 10, 2, or
// 16.
func (l *Lexer) readDecimal() token.Token {

	str := ""

	// We usually just accept digits, plus the negative unary marker.
	accept := "-0123456789"

	// But if we have `0x` as a prefix we accept hexadecimal instead.
	if l.ch == '0' && l.peekChar() == 'x' {
		accept = "0x123456789abcdefABCDEF"
	}

	// If we have `0b` as a prefix we accept binary digits only.
	if l.ch == '0' && l.peekChar() == 'b' {
		accept = "b01"
	}

	// While we have a valid character append it to our
	// result and keep reading/consuming characters.
	for strings.Contains(accept, string(l.ch)) {
		str += string(l.ch)
		l.readChar()
	}

	// If we have a `-` it can only occur at the beginning
	for _, chr := range []string{"-"} {
		if strings.Contains(str, chr) {
			if !strings.HasPrefix(str, chr) {
				return token.Token{
					Type:    token.ILLEGAL,
					Literal: "'" + chr + "' may only occur at the start of the number",
				}
			}
		}
	}

	// OK convert to an integer, which we'll later turn to a string.
	//
	// We do this so we can convert 0xff -> "255", or "0b0011" to "3".
	val, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		tok := token.Token{Type: token.ILLEGAL, Literal: err.Error()}
		return tok
	}

	// Now return that number as a string.
	return token.Token{Type: token.NUMBER, Literal: fmt.Sprintf("%v", val)}
}

// read Identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isIdentifier(l.ch) {
		l.readChar()
	}
	return string(l.characters[position:l.position])
}

// skip white space
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// read string
func (l *Lexer) readString() (string, error) {
	out := ""

	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
		if l.ch == rune(0) {
			return "", errors.New("unterminated string")
		}

		//
		// Handle \n, \r, \t, \", etc.
		//
		if l.ch == '\\' {

			// Line ending with "\" + newline
			if l.peekChar() == '\n' {
				// consume the newline.
				l.readChar()
				continue
			}

			l.readChar()

			if l.ch == rune('n') {
				l.ch = '\n'
			}
			if l.ch == rune('r') {
				l.ch = '\r'
			}
			if l.ch == rune('t') {
				l.ch = '\t'
			}
			if l.ch == rune('"') {
				l.ch = '"'
			}
			if l.ch == rune('\\') {
				l.ch = '\\'
			}
		}
		out = out + string(l.ch)

	}

	return out, nil
}

// read "[ xxxx ]"
func (l *Lexer) readEval() (string, error) {
	return l.readNestedPair('[', ']')
}

// read "{ xxxx }"
func (l *Lexer) readBlock() (string, error) {
	return l.readNestedPair('{', '}')
}

// readNestedPair reads the contents between a pair of terminators.
//
// This keeps track of the number of nested blocks so they are always
// matched.
func (l *Lexer) readNestedPair(open rune, close rune) (string, error) {

	// We're at the start of a value, but we're not yet inside it.
	depth := 1
	out := ""

	for {
		l.readChar()

		// opening a new one?
		if l.ch == open {
			out = out + string(l.ch)
			depth++
			continue
		}

		// closing a pair?
		if l.ch == close {
			depth--

			// back to zero then we're at the end of this one
			if depth == 0 {
				return out, nil
			}

			continue
		}

		// end of input?  then the user did something wrong
		if l.ch == rune(0) {
			return "", fmt.Errorf("unterminated pair %c-%c depth:%d current:%s", open, close, depth, out)
		}

		out = out + string(l.ch)

	}

}

// readVariable returns a variable
// 16.
func (l *Lexer) readVariable() string {

	str := string(l.ch)

	for {
		l.readChar()

		if l.ch == rune(0) {
			return str
		}

		if l.ch == '$' || isLetter(l.ch) {
			str += string(l.ch)
		} else {
			l.rewind()
			return str
		}
	}
}

// peek ahead at the next character
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.characters) {
		return rune(0)
	}
	return l.characters[l.readPosition]
}

// determinate whether the given character is legal within an identifier or not.
//
// This is very permissive.
func isIdentifier(ch rune) bool {
	return !isWhitespace(ch) &&
		ch != rune('{') &&
		ch != rune('}') &&
		ch != rune('[') &&
		ch != rune(']') &&
		ch != rune('$') &&
		ch != rune(';') &&
		ch != rune('\n') &&
		ch != rune(0)
}

// Is the character white space?
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

// Is the given character a digit?
func isDigit(ch rune) bool {
	return rune('0') <= ch && ch <= rune('9')
}

// Is the given character a letter?
func isLetter(ch rune) bool {
	return rune('a') <= ch && ch <= rune('z')
}
