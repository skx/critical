// Package parser contains a simple parser for processing TCL-like
// programs.
//
// It assumes that the input consists of a series of commands, which
// have an arbitrary number of associated arguments.
//
// Commands are separated from each other by either ";" or "newlines".
//
// Arguments consist of strings, numbers, evaluated-objects or blocks.
//
// We consider anything surrounded by [ & ] to be an evaluated block, and
// anything surrounded by { & } to be a block.  The latter are returned
// as-is, without any interpolation.
//
// This is a little naive, and probably doesn't handle the full scope of
// input, but it should be reasonable regardless.
//
// FIXME:
//
// -  We probably want a "variable" type.
// -  We probably parse "$foo" as a variable, and also "$foo$bar".
//    If we handle the latter as two variables we'll have problems with:
//
//       set a pu
//       set b ts
//       $a$b "Hello, world"
//
//    In practice that's probably OK, but ..
//
package parser

import (
	"fmt"

	"github.com/skx/critical/lexer"
	"github.com/skx/critical/token"
)

// Command is a single command, or statement, which should be executed.
type Command struct {
	// Command is the command to be executed.
	Command token.Token

	// Arguments contains the arguments to be used for the command.
	Arguments []token.Token
}

// Parser holds the objects' state
type Parser struct {

	// The lexer we work with for parsing the program.
	lexer *lexer.Lexer
}

// New creates a new parser.
func New(input string) *Parser {
	return &Parser{lexer: lexer.New(input)}
}

// Parse parses the input into a series of commands.
func (p *Parser) Parse() ([]Command, error) {

	// Return value, the parsed set of commands
	ret := []Command{}

	// Forever
	for {

		// The command we're processing at the moment.
		var c Command

		// Read a token
		tok := p.lexer.NextToken()

		// All done?
		if tok.Type == token.EOF {
			break
		}

		if tok.Type == token.NEWLINE {
			continue
		}

		// Some kind of error?
		if tok.Type == token.ILLEGAL {
			return ret, fmt.Errorf("illegal token:%s", tok)
		}

		// Save the token as the command to be executed.
		c.Command = tok

		// Now look for arguments to the command
		tok = p.lexer.NextToken()

		// Commands are terminated by either:
		//
		//   semi-colon (";")
		//   newline
		//   end of input
		for tok.Type != token.SEMICOLON &&
			tok.Type != token.NEWLINE &&
			tok.Type != token.EOF {

			// Add the token as an argument
			c.Arguments = append(c.Arguments, tok)

			// Read the next token
			tok = p.lexer.NextToken()
		}

		// Append the parsed command to our list,
		// and start again processing the next command.
		//
		// When we hit EOF/error we'll break out of this loop
		ret = append(ret, c)
	}

	// We hit EOF, so return the list of processed commands.
	return ret, nil
}
