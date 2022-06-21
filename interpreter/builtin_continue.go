package interpreter

import (
	"errors"
	"fmt"
)

var (
	errContinue = errors.New("CONTINUE outside a loop")
)

// continueFn is the golang implementation of the TCL `continue` function.
func continueFn(i *Interpreter, args []string) (string, error) {

	if len(args) != 0 {
		return "", fmt.Errorf("continue takes zero arguments")
	}

	return "", errContinue
}
