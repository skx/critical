package interpreter

// comment is a function which ignores comments "// xx" or "# xxxx".
func comment(i *Interpreter, args []string) (string, error) {
	return "", nil
}
