// Package interpreter is the core of our application, and is responsible
// for taking some source-code, parsing it, and evaluating it.
package interpreter

import (
	"fmt"
	"os"
	"regexp"

	"github.com/skx/critical/environment"
	"github.com/skx/critical/parser"
	"github.com/skx/critical/token"
)

// HostFunctionSignature is the signature of a function implemented in golang, which
// is exported to the TCL scripting environment.
type HostFunctionSignature func(i *Interpreter, args []string) (string, error)

// HostFunction represents a built-in function available to the TCL environment
// which is implemented in golang.
type HostFunction struct {

	// function is the golang function to handle the call
	function HostFunctionSignature
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

	// builtins contains pointers to the golang implementations of
	// the TCL functions.
	builtins map[string]HostFunction

	// functions contain user-defined functions, written in TCL.
	functions map[string]UserFunction

	// program is the parsed program.
	program []parser.Command

	// environment holds any variable-references the user has defined.
	//
	// Note that functions the user defines are not stored here, they
	// live in the `functions` map.
	environment *environment.Environment
}

// New creates a new object to interpret.
func New(source string) (*Interpreter, error) {

	// Create the object we'll return
	i := &Interpreter{
		builtins:    make(map[string]HostFunction),
		environment: environment.New(),
		functions:   make(map[string]UserFunction),
	}

	// parser is the object we use to transform the source into
	// a program we can evaluate.
	parser := parser.New(source)

	// Parse the program to find any obvious errors immediately.
	var err error
	i.program, err = parser.Parse()
	if err != nil {
		return nil, err
	}

	// Bind the expected primitives
	i.RegisterBuiltin("append", appendFn)
	i.RegisterBuiltin("break", breakFn)
	i.RegisterBuiltin("continue", continueFn)
	i.RegisterBuiltin("decr", decr)
	i.RegisterBuiltin("env", env)
	i.RegisterBuiltin("eval", evalFn)
	i.RegisterBuiltin("exit", exitFn)
	i.RegisterBuiltin("expr", expr)
	i.RegisterBuiltin("for", forFn)
	i.RegisterBuiltin("if", ifFn)
	i.RegisterBuiltin("incr", incr)
	i.RegisterBuiltin("proc", proc)
	i.RegisterBuiltin("puts", puts)
	i.RegisterBuiltin("regexp", regexpFn)
	i.RegisterBuiltin("return", returnFn)
	i.RegisterBuiltin("set", set)
	i.RegisterBuiltin("while", while)

	return i, nil
}

// Evaluate parses the program source, and executes the program.
func (i *Interpreter) Evaluate() (string, error) {

	// Output of the evaluation is the output received from the
	// last statement which was executed.
	out := ""

	// Holder for any error we might receive
	var err error

	// For each parsed command, evaluate it
	for _, cmd := range i.program {

		// The name of the command we're going to run
		name := ""

		// The name might require expansion, so handle that first.
		switch cmd.Command.Type {
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

			// If the function returned a value then use that.
			if e == ErrReturn {
				return out, e
			}

			//
			// `break` and `continue` errors are handled specially
			// within the handlers for `if` and `while`.
			//
			// Here we just return them, and they'll do the
			// right thing.
			//
			// The same thing applies more generally to the
			// exit handler
			//
			if e == errBreak || e == errContinue || e == ErrExit {
				return out, e
			}

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

			// If the function returned a value then use that.
			if e == ErrReturn {
				continue
			}

			// Exit inside a proc.
			if e == ErrExit {
				return out, e
			}

			// Now we've restored the environment we can
			// handle the error-detection
			if e != nil {
				return "", e
			}

			continue
		}

		// At this point we've been given a "command" which
		// doesn't exist as a function - either in golang, or
		// user-defined.
		//
		// If the input was a literal string, number, or variable
		// to be expanded then we set our return value to that.
		//
		if cmd.Command.Type == token.STRING {
			out = cmd.Command.Literal
			continue
		}
		if cmd.Command.Type == token.NUMBER {
			out = cmd.Command.Literal
			continue
		}
		if cmd.Command.Type == token.VARIABLE {
			// The actual literal will have been updated
			out = i.expandString(cmd.Command.Literal)
			continue
		}

		//
		// Otherwise we just return an error.
		//
		return "", fmt.Errorf("unknown command '%s':%v", name, cmd)

	}
	return out, err
}

// expandString converts "$foo $bar" into "$ENV{'FOO'} $ENV{'BAR'}".
func (i *Interpreter) expandString(str string) string {

	// Lookup any variables inside our environment.
	//
	// If they're not found then return ""
	//
	mapper := func(placeholderName string) string {
		if placeholderName == "$" {
			return "$"
		}
		val, ok := i.environment.Get(placeholderName)
		if ok {
			return val
		}
		return ""
	}

	// Replace the string, using the mapper.
	return os.Expand(str, mapper)
}

// Eval handles sub-expressions, this is horrid.
// TODO / FIXME / HACK / XXX
func (i *Interpreter) Eval(str string) (string, error) {

	// sub-evaluator.  horrid
	tmp, er := New(str)
	if er != nil {
		return "", er
	}

	// Hacky way to ensure that any updates to the variables
	// affect this environment too, not just the child one.
	tmp.environment = i.environment

	// Hacky way to ensure the child environment has the same
	// functions as we do.
	tmp.functions = i.functions
	tmp.builtins = i.builtins

	// run the script
	out, err := tmp.Evaluate()

	if err == ErrReturn || err == ErrExit {
		return out, err
	}
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

		// And continue going until we find no more matches.
		out = r.FindStringSubmatch(str)
	}

	return str
}

// RegisterBuiltin registers a builtin function.
func (i *Interpreter) RegisterBuiltin(name string, fn HostFunctionSignature) {
	i.builtins[name] = HostFunction{function: fn}
}
