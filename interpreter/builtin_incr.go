package interpreter

import (
	"fmt"
	"strconv"
)

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
	increase := 1.0
	if len(args) == 2 {
		var err error
		increase, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return "", err
		}
	}

	orig, _ := strconv.ParseFloat(cur, 64)
	orig += increase

	// an integer, really?
	if orig == float64(int(orig)) {
		i.environment.Set(name, fmt.Sprintf("%d", int(orig)))
		return fmt.Sprintf("%d", int(orig)), nil
	}

	// A floating-point number
	i.environment.Set(name, fmt.Sprintf("%f", orig))
	return fmt.Sprintf("%f", orig), nil
}
