package interpreter

import "fmt"

// puts is the golang implementation of the TCL `puts` function.
func puts(i *Interpreter, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("puts only accepts one argument, got %d", len(args))
	}

	fmt.Printf("%s\n", args[0])
	return args[0], nil
}
