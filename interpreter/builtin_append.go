package interpreter

import "fmt"

// appendFn is the golang implementation of the TCL `append` function.
func appendFn(i *Interpreter, args []string) (string, error) {

	if len(args) < 1 {
		return "", fmt.Errorf("append requires at least one argument")
	}

	// Get the value of the variable
	val, _ := i.environment.Get(args[0])

	// Append all the arguments to it - skipping the first
	for i, d := range args {
		if i == 0 {
			continue
		}

		val += d
	}

	// Update the value, and also return it
	i.environment.Set(args[0], val)
	return val, nil
}
