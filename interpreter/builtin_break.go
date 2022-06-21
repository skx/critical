package interpreter

import (
	"errors"
	"fmt"
)

var (
	errBreak = errors.New("BREAK outside a loop")
)

// breakFn is the golang implementation of the TCL `break` function.
func breakFn(i *Interpreter, args []string) (string, error) {

	if len(args) != 0 {
		return "", fmt.Errorf("break takes zero arguments")
	}

	return "", errBreak
}
