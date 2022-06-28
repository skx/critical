package interpreter

import "fmt"

// forFn is the golang implementation of the TCL `for` function.
func forFn(i *Interpreter, args []string) (string, error) {

	// for {setup} {test} {incr} {body}
	// Test arguments
	if len(args) != 4 {
		return "", fmt.Errorf("for requires four arguments.  Got %d", len(args))
	}

	pro := args[0]
	mid := args[1]
	epi := args[2]
	bdy := args[3]

	// return value
	var out string

	// temporary, used for evaling the test
	var tmp string

	// error holder
	var err error

	// Evaluate the prologue
	_, err = i.Eval(pro)
	if err != nil {
		return "", err
	}

	// Now the condition
	for {

		// middle part
		tmp, err = i.Eval(mid)
		if err != nil {
			return "", err
		}

		// if the condition fails then we're out.
		if tmp == "" || tmp == "0" {
			return out, nil
		}

		// condition was true, so run the loop body
		out, err = i.Eval(bdy)

		//
		// We might have a synthetic-error for the
		// control words "break" & "continue".
		//
		if err == errBreak {

			// GOTO Considered useful
			goto outside

		} else if err == errContinue {

			// Nop

		} else if err != nil {

			// Any other error is fatal.
			return "", err
		}

		// now the post-condition
		_, err = i.Eval(epi)
		if err != nil {
			return "", err
		}
	}

outside:
	// Return the last statement from within the body
	return out, nil

}
