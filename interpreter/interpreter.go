// Package interpreter is the core of our application, and is responsible
// for taking some source-code, parsing it, and evaluating it.
package interpreter

import (
	"fmt"
	"regexp"

	"github.com/skx/critical/environment"
	"github.com/skx/critical/parser"
	"github.com/skx/critical/token"
)

// BuiltIn represents a built-in function
type BuiltIn struct {
	// function is the golang function to handle the call
	function func(i *Interpreter, args []string) (string, error)
}

// Interpreter holds the interpreters state.
type Interpreter struct {

	// Source is the program we're going to execute
	source string

	// parser is the object we use to transform the source into
	// a program we can evaluate
	parser *parser.Parser

	// builtins contains pointers to the golang implementations of
	// the TCL functions.
	builtins map[string]BuiltIn

	// environment contains any variables the user has defined
	environment *environment.Environment
}

// New creates a new object to interpret.
func New(source string) *Interpreter {

	// Create the object we'll return
	i := &Interpreter{
		builtins:    make(map[string]BuiltIn),
		environment: environment.New(),
		parser:      parser.New(source),
		source:      source,
	}

	// Bind the expected primitives

	// These are internal functions that aren't real
	i.builtins["#"] = BuiltIn{function: comment}
	i.builtins["//"] = BuiltIn{function: comment}
	i.builtins["\\n"] = BuiltIn{function: comment}

	// These are real primitives
	i.builtins["decr"] = BuiltIn{function: decr}
	i.builtins["expr"] = BuiltIn{function: expr}
	i.builtins["if"] = BuiltIn{function: iff}
	i.builtins["incr"] = BuiltIn{function: incr}
	i.builtins["puts"] = BuiltIn{function: puts}
	i.builtins["set"] = BuiltIn{function: set}
	i.builtins["while"] = BuiltIn{function: while}

	return i
}

// Evaluate parses the program source, and executes the program.
func (i *Interpreter) Evaluate(allowEmpty bool) (string, error) {

	program, err := i.parser.Parse()
	if err != nil {
		return "", err
	}

	// Output of the evaluation is the last output
	out := ""

	// For each parsed command, evaluate it
	for _, cmd := range program {

		// The name of the command we're going to run
		name := ""

		// The name might require expansion, so handle that first
		switch cmd.Command.Type {
		case token.NEWLINE:
			// Yes, this is crazy
			name = "\\n"
		case token.STRING:
			name = cmd.Command.Literal
		case token.NUMBER:
			name = cmd.Command.Literal
		case token.IDENT:
			name = cmd.Command.Literal
		case token.VARIABLE:
			name = i.expandString(cmd.Command.Literal)
		default:
			return "", fmt.Errorf("unknown command type %v", cmd.Command)
		}

		// We need to expand the arguments to the command, which
		// will convert them into the appropriate arguments.
		var args []string

		// This will need some love in the future.
		for _, arg := range cmd.Arguments {

			// If the token isn't a quoted form we've got
			// expand it.
			if arg.Type != token.BLOCK {

				// Expand "$a" -> "$ENV{a}"
				expand := i.expandString(arg.Literal)

				// Now expand "[ .. x .. ]" to include
				// having that there.
				expand = i.expandEval(expand)

				args = append(args, expand)
			} else {
				args = append(args, arg.Literal)
			}
		}

		fn, ok := i.builtins[name]
		if !ok {

			// OK so the thing wasn't
			if allowEmpty {
				return name, nil
			}
			return "", fmt.Errorf("unknown command '%s'", name)
		}

		// Call the function, and if it errors then abort
		var e error
		out, e = fn.function(i, args)
		if e != nil {
			return "", fmt.Errorf("error invoking %s: %s", cmd.Command.Literal, e)
		}

	}
	return out, err
}

func (i *Interpreter) expandString(str string) string {
	ret := ""

	idx := 0

	for idx < len(str) {

		if str[idx] == '$' {

			// Skip past the dollar
			idx++

			// We build up the name of the variable to
			// expand
			variable := ""

			// While we've not walked off the end of our
			// string, and we've got a "letter" then we
			// can update our variable name.
			for idx < len(str) && isLetter(str[idx]) {
				variable += string(str[idx])
				idx++
			}

			// OK append the variable value to the string
			val, _ := i.environment.Get(variable)
			ret += val
		} else {

			// Just append the string
			ret += str[idx : idx+1]
			idx++
		}
	}

	return ret
}

// Eval handles sub-expressions, this is horrid
// TODO / FIXME / HACK / XXX
func (i *Interpreter) Eval(str string) (string, error) {

	// sub-evaluator.  horrid
	tmp := New(str)

	// Hacky way to ensure that any updates to the variables
	// affect this environment too, not just the child one.
	tmp.environment = i.environment

	// run the script
	out, err := tmp.Evaluate(true)
	if err != nil {
		return "", err
	}

	return out, nil
}

// expandEval handles the expansion of "[ FOO ]" blocks.
//
// It now handles nested values, such that this works:
//
//    puts [ expr 1 + [ expr 2 + 3 ] ]
//
func (i *Interpreter) expandEval(str string) string {

	r := regexp.MustCompile(`^(.*)\[([^\]]+)\](.*)$`)

	out := r.FindStringSubmatch(str)
	for len(out) > 1 {

		// The pieces of the match
		before := out[1]
		match := out[2]
		after := out[3]

		// Evaluate the middle.
		eval, _ := i.Eval(match)

		// Now update our string.
		str = before + eval + after

		out = r.FindStringSubmatch(str)
	}

	return str
}

//
// Built-In functions here ..
//

func isLetter(ch byte) bool {
	return ((ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z'))
}
