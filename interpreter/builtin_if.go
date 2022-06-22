package interpreter

import "fmt"

// ifFn is the golang implementation of the TCL `if` function.
func ifFn(i *Interpreter, args []string) (string, error) {

	// Test arguments
	if len(args) != 2 && len(args) != 4 {
		return "", fmt.Errorf("if accepts two arguments, or four.  Got %d", len(args))
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
	}

	// If we have a "false branch", then execute it
	if fail != "" {
		return i.Eval(fail)
	}

	return "", nil
}
