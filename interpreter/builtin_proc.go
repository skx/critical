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

	// args - split by space
	argsIn := strings.Split(args[1], " ")

	// But only collect the non-empty ones
	argsOut := []string{}

	for _, arg := range argsIn {
		if strings.TrimSpace(arg) != "" {
			argsOut = append(argsOut, arg)
		}
	}

	// body
	body := args[2]

	// Save the function
	i.functions[name] = UserFunction{
		Args: argsOut,
		Body: body,
	}

	return "", nil
}
