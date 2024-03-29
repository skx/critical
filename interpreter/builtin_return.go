package interpreter

import (
	"errors"
	"fmt"
)

var (
	// ErrReturn will be used to handle return-values from functions.
	//
	// It should be handled and expected by callers.
	ErrReturn = errors.New("RETURN")
)

// returnFn is the golang implementation of the TCL `return` function.
func returnFn(i *Interpreter, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("return takes one argument")
	}

	return args[0], ErrReturn
}
