package interpreter

import (
	"fmt"
	"strconv"
)

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
	decrease := 1.0
	if len(args) == 2 {
		var err error
		decrease, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return "", err
		}
	}

	orig, _ := strconv.ParseFloat(cur, 64)
	orig -= decrease

	// an integer, really?
	if orig == float64(int(orig)) {
		i.environment.Set(name, fmt.Sprintf("%d", int(orig)))
		return fmt.Sprintf("%d", int(orig)), nil
	}

	i.environment.Set(name, fmt.Sprintf("%f", orig))
	return fmt.Sprintf("%f", orig), nil

}
