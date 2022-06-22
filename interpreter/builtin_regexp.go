package interpreter

import (
	"fmt"
	"regexp"
)

// regexpFn is the golang implementation of the TCL `regexp` function.
func regexpFn(i *Interpreter, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("regexp requires two arguments, got %d", len(args))
	}

	r, err := regexp.Compile(args[0])
	if err != nil {
		return "", err
	}

	if r.MatchString(args[1]) {
		return "1", nil
	}
	return "0", nil
}
