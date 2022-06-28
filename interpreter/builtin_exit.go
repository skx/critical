package interpreter

import (
	"errors"
	"fmt"
)

var (
	// ErrExit is the "error" code a script will terminate with if
	// it finishes execution with an `exit` statement.
	ErrExit = errors.New("EXIT")
)

// exitFn is the golang implementation of the TCL `exit` function.
func exitFn(i *Interpreter, args []string) (string, error) {

	// Ensure we have only one argument
	if len(args) != 1 {
		return "", fmt.Errorf("exit only accepts one argument, got %d", len(args))
	}

	return args[0], ErrExit
}
