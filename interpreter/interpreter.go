// Package interpreter is the core of our application, and is responsible
// for taking some source-code, parsing it, and evaluating it.
package interpreter

import (
	"fmt"
	"strconv"

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

func (i *Interpreter) expandEval(str string) string {
	ret := ""
	idx := 0

	for idx < len(str) {
		c := str[idx]

		// OK we have a nested thing.
		if c == '[' {
			tmp := ""
			closed := false
			idx++
			for idx < len(str) && !closed {
				c = str[idx]
				idx++
				if c == ']' {
					out, _ := i.Eval(tmp)
					ret += out
					closed = true

				} else {
					tmp += string(c)
				}
			}

		} else {
			ret += string(c)
		}
		idx++
	}
	return ret
}

//
// Built-In functions here ..
//

// comment is a function which ignores comments "// xx" or "# xxxx".
func comment(i *Interpreter, args []string) (string, error) {
	return "", nil
}

// decr is the golang implementation of the TCL `decr` function.
func decr(i *Interpreter, args []string) (string, error) {

	if len(args) != 1 && len(args) != 2 {
		return "", fmt.Errorf("decr takes one or two arguments")
	}

	// Name of variable we're decreasing
	name := args[0]

	// Get the current value of the variable
	// if not found the value is zero
	cur, ok := i.environment.Get(name)
	if !ok {
		cur = "0"
	}

	// How much to decrease by?
	decrease := 1
	if len(args) == 2 {
		var err error
		decrease, err = strconv.Atoi(args[1])
		if err != nil {
			return "", nil
		}
	}

	orig, _ := strconv.Atoi(cur)
	orig -= decrease
	i.environment.Set(name, fmt.Sprintf("%d", orig))

	return fmt.Sprintf("%d", orig), nil

}

// expr is the golang implementation of the TCL `expr` function.
func expr(i *Interpreter, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("expr requires three arguments, got %d", len(args))
	}

	aV, eA := strconv.Atoi(args[0])
	if eA != nil {
		return "", eA
	}
	op := args[1]
	bV, eB := strconv.Atoi(args[2])
	if eB != nil {
		return "", eB
	}

	switch op {
	case "+":
		return (fmt.Sprintf("%f", float64(aV+bV))), nil
	case "-":
		return (fmt.Sprintf("%f", float64(aV-bV))), nil
	case "*":
		return (fmt.Sprintf("%f", float64(aV*bV))), nil
	case "/":
		return (fmt.Sprintf("%f", float64(aV/bV))), nil
	case "<":
		if aV < bV {
			return "1", nil
		} else {
			return "0", nil
		}
	case "<=":
		if aV <= bV {
			return "1", nil
		} else {
			return "0", nil
		}
	case ">":
		if aV > bV {
			return "1", nil
		} else {
			return "0", nil
		}
	case ">=":
		if aV >= bV {
			return "1", nil
		} else {
			return "0", nil
		}
	}
	return "", fmt.Errorf("unknown operation %s %s %s", args[0], op, args[2])
}

// iff is the golang implementation of the TCL `if` function.
func iff(i *Interpreter, args []string) (string, error) {

	// Test arguments
	if len(args) != 2 && len(args) != 4 {
		return "", fmt.Errorf("if accepts three arguments, or five.  Got %d", len(args))
	}

	cond := args[0]
	pass := args[1]
	fail := ""

	if len(args) == 4 {
		fail = args[3]
	}

	//
	// evaluate the condition.
	//
	// if true
	//   eval(pass)
	// else
	//   if fail != ""
	//     eval(fail)
	//

	out, err := i.Eval(cond)
	if err != nil {
		return "", err
	}

	// A non-false result means we run the "true branch"
	if out != "" && out != "0" {
		return i.Eval(pass)
	} else {
		// If we have a "false branch", then execute it
		if fail != "" {
			return i.Eval(fail)
		}
	}

	return "", nil
}

// incr is the golang implementation of the TCL `incr` function.
func incr(i *Interpreter, args []string) (string, error) {

	if len(args) != 1 && len(args) != 2 {
		return "", fmt.Errorf("incr takes one or two arguments")
	}

	// Name of variable we're increasing
	name := args[0]

	// Get the current value of the variable
	// if not found the value is zero
	cur, ok := i.environment.Get(name)
	if !ok {
		cur = "0"
	}

	// How much to increase by?
	increase := 1
	if len(args) == 2 {
		var err error
		increase, err = strconv.Atoi(args[1])
		if err != nil {
			return "", nil
		}
	}

	orig, _ := strconv.Atoi(cur)
	orig += increase
	i.environment.Set(name, fmt.Sprintf("%d", orig))

	return fmt.Sprintf("%d", orig), nil
}

// puts is the golang implementation of the TCL `puts` function.
func puts(i *Interpreter, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("puts only accepts one argument, got %d", len(args))
	}

	fmt.Printf("%s\n", args[0])
	return args[0], nil
}

// set is the golang implementation of the TCL `set` function.
func set(i *Interpreter, args []string) (string, error) {
	if len(args) != 1 && len(args) != 2 {
		return "", fmt.Errorf("set accepts one or two arguments, got %d", len(args))
	}

	// name of the value we're getting/setting
	name := args[0]

	// If we have a value, then set it and return it.
	if len(args) == 2 {
		value := args[1]
		i.environment.Set(name, value)
		return value, nil
	}

	// otherwise return the current value
	cur, _ := i.environment.Get(name)
	return cur, nil
}

// while is the golang implementation of the TCL `while` function.
func while(i *Interpreter, args []string) (string, error) {

	// Test arguments
	if len(args) != 2 {
		return "", fmt.Errorf("while accepts only two arguments.  Got %d", len(args))
	}

	cond := args[0]
	body := args[1]

	out := ""
	var err error
	tmp := ""

	// Run the conditional once
	tmp, err = i.Eval(cond)
	if err != nil {
		return "", err
	}

	// Run the body, and repeat until the conditional fails
	for tmp != "0" && tmp != "" {

		// run the body
		out, err = i.Eval(body)

		// repeat the
		tmp, err = i.Eval(cond)
		if err != nil {
			return "", err
		}
	}

	// Return the last statement from within the body
	return out, nil
}

func isLetter(ch byte) bool {
	return ((ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z'))
}
