package interpreter

import (
	"fmt"
	"os"
	"strconv"
)

// exitFn is the golang implementation of the TCL `exit` function.
func exitFn(i *Interpreter, args []string) (string, error) {

	// Ensure we have only one argument
	if len(args) != 1 {
		return "", fmt.Errorf("exit only accepts one argument, got %d", len(args))
	}

	// Convert to an integer
	var num int
	var err error

	num, err = strconv.Atoi(args[0])
	if err != nil {

		// error?
		num = 1
	}

	os.Exit(num)

	return "unreached", nil
}
