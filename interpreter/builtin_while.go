package interpreter

import "fmt"

// while is the golang implementation of the TCL `while` function.
func while(i *Interpreter, args []string) (string, error) {

	// Test arguments
	if len(args) != 2 {
		return "", fmt.Errorf("while accepts only two arguments.  Got %d", len(args))
	}

	cond := args[0]
	body := args[1]

	out := ""
	var err error
	tmp := ""

	// Run the conditional once
	tmp, err = i.Eval(cond)
	if err != nil {
		return "", err
	}

	// Run the body, and repeat until the conditional fails
	for tmp != "0" && tmp != "" {

		// run the body
		out, err = i.Eval(body)

		// We might have BREAK or CONTINUE within the loop.
		if err == errBreak {

			// GOTO Considered useful
			goto outside

		} else if err == errContinue {

			// Nop

		} else if err == ErrExit {

			// Exit
			return out, err

		} else if err != nil {

			// Another, unexpected error
			return out, err
		}

		// repeat the conditional-test ahead of repeating the body
		tmp, err = i.Eval(cond)
		if err != nil {
			return "", err
		}
	}
outside:
	// Return the last statement from within the body
	return out, nil
}
