package interpreter

import "fmt"

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
