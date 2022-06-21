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

// HostFunction represents a built-in function available to the TCL environment
// which is implemented in golang.
type HostFunction struct {

	// function is the golang function to handle the call
	function func(i *Interpreter, args []string) (string, error)
}

// UserFunction represents a function which has been defined by the user,
// within the TCL environment (via the use of 'proc').
type UserFunction struct {

	// Args contains the list of parameters
	Args []string

	// Body contains the function body
	Body string
}

// Interpreter holds the interpreters state.
type Interpreter struct {

	// parser is the object we use to transform the source into
	// a program we can evaluate.
	parser *parser.Parser

	// builtins contains pointers to the golang implementations of
	// the TCL functions.
	builtins map[string]HostFunction

	// functions contain user-defined functions, written in TCL.
	functions map[string]UserFunction

	// environment holds any variable-references the user has defined.
	//
	// Note that functions the user defines are not stored here, they
	// live in the `functions` map.
	environment *environment.Environment
}

// New creates a new object to interpret.
func New(source string) *Interpreter {

	// Create the object we'll return
	i := &Interpreter{
		builtins:    make(map[string]HostFunction),
		environment: environment.New(),
		functions:   make(map[string]UserFunction),
		parser:      parser.New(source),
	}

	// Bind the expected primitives

	// These are internal functions that aren't real
	i.builtins["#"] = HostFunction{function: comment}
	i.builtins["//"] = HostFunction{function: comment}
	i.builtins["\\n"] = HostFunction{function: comment}

	// These are real primitives
	i.builtins["decr"] = HostFunction{function: decr}
	i.builtins["expr"] = HostFunction{function: expr}
	i.builtins["if"] = HostFunction{function: iff}
	i.builtins["incr"] = HostFunction{function: incr}
	i.builtins["puts"] = HostFunction{function: puts}
	i.builtins["proc"] = HostFunction{function: proc}
	i.builtins["set"] = HostFunction{function: set}
	i.builtins["while"] = HostFunction{function: while}

	return i
}

// Evaluate parses the program source, and executes the program.
func (i *Interpreter) Evaluate() (string, error) {

	// Parse the program, if there were errors bail immediately.
	program, err := i.parser.Parse()
	if err != nil {
		return "", err
	}

	// Output of the evaluation is the output received from the
	// last statement which was executed.
	out := ""

	// For each parsed command, evaluate it
	for _, cmd := range program {

		// The name of the command we're going to run
		name := ""

		// The name might require expansion, so handle that first.
		switch cmd.Command.Type {
		case token.NEWLINE:
			// Yes, this is crazy.
			//
			// It might no longer be necessary but in the past
			// my parser would return newline-objects if there
			// were multiple blank lines.
			//
			// Rather than wading through to fix that I just
			// pretend there's a command called "\n" which then
			// throws away all arguments/input.
			//
			// This is also used to handle comments, prefixed with
			// either "#" or "//".
			//
			// It's almost sane, but ..
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

		// We need to expand the arguments to the command, so here
		// is the place to store those converted args, before we
		// pass them to the handler.
		var args []string

		// For each argument
		for _, arg := range cmd.Arguments {

			// If the token isn't a quoted form we've got
			// expand it.
			if arg.Type != token.BLOCK {

				// Expand "$a" -> "$ENV{a}"
				expand := i.expandString(arg.Literal)

				// Now expand "[ .. x .. ]" to include
				// having that there.
				expand = i.expandEval(expand)

				// Save the argument away.
				args = append(args, expand)
			} else {
				// This is a quoted-block, just append literally
				args = append(args, arg.Literal)
			}
		}

		// Is the function a built-in implemented in golang?
		fn, ok := i.builtins[name]
		if ok {

			// Call the function, and if it errors then abort
			var e error
			out, e = fn.function(i, args)
			if e != nil {
				return "", fmt.Errorf("error invoking %s: %s", cmd.Command.Literal, e)
			}
			continue
		}

		// Is the function a user-written function in TCL?
		userFN, ok2 := i.functions[name]

		if ok2 {
			var e error

			if len(args) != len(userFN.Args) {
				return "", fmt.Errorf("function argument mismatch, %s takes %d arguments, %d supplied", name, len(userFN.Args), len(args))
			}

			// Save old environment
			oldE := i.environment

			// Create a new environment
			newE := environment.NewEnclosedEnvironment(oldE)

			// Make the environment live
			i.environment = newE

			// Set the environment variables for the proc
			// arguments.
			for idx, arg := range userFN.Args {
				i.environment.Set(arg, args[idx])
			}

			out, e = i.Eval(userFN.Body)

			// Restore the old environment, now the function
			// is over.
			i.environment = oldE

			// Now we've restored the environment we can
			// handle the error-detection
			if e != nil {
				fmt.Printf("Error calling user function:%s\n", e)
				return "", e
			}

			continue
		}

		// At this point we've been given a "command" which
		// doesn't exist as a function - either in golang, or
		// user-defined.
		//
		// If the input was a literal string then we'll return
		// that
		if cmd.Command.Type == token.STRING {
			return cmd.Command.Literal, nil
		}
		if cmd.Command.Type == token.NUMBER {
			return cmd.Command.Literal, nil
		}

		//
		// Otherwise we just return an error.
		//
		return "", fmt.Errorf("unknown command '%s'", name)

	}
	return out, err
}

// expandString converts "$foo $bar" into "$ENV{'FOO'} $ENV{'BAR'}".
//
// It's a little horrid.
//
// TODO / FIXME / HACKY
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

// Eval handles sub-expressions, this is horrid.
// TODO / FIXME / HACK / XXX
func (i *Interpreter) Eval(str string) (string, error) {

	// sub-evaluator.  horrid
	tmp := New(str)

	// Hacky way to ensure that any updates to the variables
	// affect this environment too, not just the child one.
	tmp.environment = i.environment

	// Hacky way to ensure the child environment has the same
	// functions as we do.
	tmp.functions = i.functions

	// run the script
	out, err := tmp.Evaluate()
	if err != nil {
		return "", err
	}

	return out, nil
}

// expandEval handles the expansion of "[ FOO ]" blocks.
//
// It correctly handles nested values, such that this works:
//
//    puts [ expr 1 + [ expr 2 + 3 ] ]
//
// The downside is we use a regular expression to handle the nested
// processing, and that means we've got three problems..
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

// isLetter should be removed, or improved.  It is used for the variable
// expansion only, via `expandString`.
//
// TODO / FIXME / REMOVEME
func isLetter(ch byte) bool {
	return ((ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z'))
}
