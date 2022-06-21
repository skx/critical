package interpreter

import "fmt"

// comment is a function which ignores comments "// xx" or "# xxxx".
func comment(i *Interpreter, args []string) (string, error) {

	if len(args) != 0 {
		return "", fmt.Errorf("comment takes zero arguments")
	}

	return "", nil
}
