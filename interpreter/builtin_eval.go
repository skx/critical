package interpreter

import "fmt"

// evalFn is the golang implementation of the TCL `eval` function.
func evalFn(i *Interpreter, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("eval only accepts one argument, got %d", len(args))
	}

	return i.Eval(args[0])
}
