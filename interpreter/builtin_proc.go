package interpreter

import (
	"fmt"
	"strings"
)

// proc is the golang implemention of the TCL `proc` function
func proc(i *Interpreter, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("proc only accepts three argument, got %d", len(args))
	}

	// name
	name := args[0]

	// args
	parm := args[1]
	a := strings.Split(parm, " ")

	// body
	body := args[2]

	// Save the function
	i.functions[name] = UserFunction{
		Args: a,
		Body: body,
	}

	return "", nil
}
