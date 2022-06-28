package interpreter

import (
	"fmt"
	"strconv"
)

var (
	// ops is a map for the function to invoke for various numeric
	// operations
	//
	// We use a map to avoid the complexity-cost of using a switch
	// statement.
	ops map[string]func(a float64, b float64) (string, error)
)

func init() {

	// Create the map
	ops = make(map[string]func(a float64, b float64) (string, error))

	// populate it with basic operations
	ops["+"] = plusFn
	ops["-"] = minusFn
	ops["/"] = divideFn
	ops["*"] = multiplyFn
	ops["%"] = modFn

	// comparisons
	ops["<"] = lessFn
	ops["<="] = lessEqualFn
	ops[">"] = greaterFn
	ops[">="] = greaterEqualFn

	// equality
	ops["=="] = eqFn
	ops["!="] = neFn
}

// expr is the golang implementation of the TCL `expr` function.
func expr(i *Interpreter, args []string) (string, error) {

	// Test argument count
	if len(args) != 3 {
		return "", fmt.Errorf("expr requires three arguments, got %d", len(args))
	}

	//
	// Parse the two arguments as numbers
	//
	// However note that the operations "ne" and "eq" are string-based
	// so if we fail to parse the params as float64s with those operations
	// things are OK
	//
	aV, eA := strconv.ParseFloat(args[0], 64)
	if eA != nil && (args[1] != "ne" && args[1] != "eq") {
		return "", eA
	}
	op := args[1]
	bV, eB := strconv.ParseFloat(args[2], 64)
	if eB != nil && (args[1] != "ne" && args[1] != "eq") {
		return "", eB
	}

	//
	// Handle the string values specially
	//
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
	}

	//
	// Lookup the operation-function to invoke in the map
	//
	opFn, ok := ops[op]
	if !ok {

		//
		// If we didn't find one that means we have an unknown
		// operation.
		//
		return "", fmt.Errorf("unknown operation %s %s %s", args[0], op, args[2])
	}

	//
	// Return whatever the operation returns.
	return opFn(aV, bV)
}

func plusFn(a float64, b float64) (string, error) {

	x := a + b

	// an integer, really?
	if x == float64(int(x)) {
		return fmt.Sprintf("%d", int(x)), nil
	}

	return (fmt.Sprintf("%f", x)), nil
}

func minusFn(a float64, b float64) (string, error) {
	x := a - b
	// an integer, really?
	if x == float64(int(x)) {
		return fmt.Sprintf("%d", int(x)), nil
	}

	return (fmt.Sprintf("%f", x)), nil
}

func multiplyFn(a float64, b float64) (string, error) {
	x := a * b
	// an integer, really?
	if x == float64(int(x)) {
		return fmt.Sprintf("%d", int(x)), nil
	}

	return (fmt.Sprintf("%f", x)), nil
}

func divideFn(a float64, b float64) (string, error) {
	x := a / b
	// an integer, really?
	if x == float64(int(x)) {
		return fmt.Sprintf("%d", int(x)), nil
	}

	return (fmt.Sprintf("%f", x)), nil
}

func modFn(a float64, b float64) (string, error) {

	return (fmt.Sprintf("%d", int(a)%int(b))), nil
}

func lessFn(a float64, b float64) (string, error) {
	if a < b {
		return "1", nil
	}
	return "0", nil
}

func lessEqualFn(a float64, b float64) (string, error) {
	if a <= b {
		return "1", nil
	}
	return "0", nil
}

func greaterFn(a float64, b float64) (string, error) {
	if a > b {
		return "1", nil
	}
	return "0", nil
}
func greaterEqualFn(a float64, b float64) (string, error) {
	if a >= b {
		return "1", nil
	}
	return "0", nil
}

func eqFn(a float64, b float64) (string, error) {
	if a == b {
		return "1", nil
	}
	return "0", nil
}

func neFn(a float64, b float64) (string, error) {
	if a != b {
		return "1", nil
	}
	return "0", nil
}
