package interpreter

import (
	"fmt"
	"os"
)

// env is the golang implementation of the TCL `env` function.
func env(i *Interpreter, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("env only accepts one argument, got %d", len(args))
	}

	return os.Getenv(args[0]), nil
}
