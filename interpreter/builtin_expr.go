package interpreter

import (
	"fmt"
	"strconv"
)

// expr is the golang implementation of the TCL `expr` function.
func expr(i *Interpreter, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("expr requires three arguments, got %d", len(args))
	}

	aV, eA := strconv.Atoi(args[0])
	if eA != nil && (args[1] != "ne" && args[1] != "eq") {
		return "", eA
	}
	op := args[1]
	bV, eB := strconv.Atoi(args[2])
	if eB != nil && (args[1] != "ne" && args[1] != "eq") {
		return "", eB
	}

	switch op {

	case "eq":
		// string equality test
		if args[0] == args[2] {
			return "1", nil
		}
		return "0", nil
	case "ne":
		// string inequality test
		if args[0] != args[2] {
			return "1", nil
		}
		return "0", nil

	case "+":
		return (fmt.Sprintf("%d", aV+bV)), nil
	case "-":
		return (fmt.Sprintf("%d", aV-bV)), nil
	case "*":
		return (fmt.Sprintf("%d", aV*bV)), nil
	case "/":
		return (fmt.Sprintf("%d", aV/bV)), nil
	case "%":
		return (fmt.Sprintf("%d", aV%bV)), nil
	case "<":
		if aV < bV {
			return "1", nil
		}
		return "0", nil
	case "==":
		if aV == bV {
			return "1", nil
		}
		return "0", nil
	case "<=":
		if aV <= bV {
			return "1", nil
		}
		return "0", nil

	case ">":
		if aV > bV {
			return "1", nil
		}
		return "0", nil

	case ">=":
		if aV >= bV {
			return "1", nil
		}
		return "0", nil
	}
	return "", fmt.Errorf("unknown operation %s %s %s", args[0], op, args[2])
}
