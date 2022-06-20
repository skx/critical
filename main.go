package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skx/critical/interpreter"
)

func main() {

	// Ensure we have a file to execute.
	if len(os.Args) < 2 {
		fmt.Printf("Usage: critical file.tcl\n")
		return
	}

	// Read the file
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("error reading file %s:%s\n", os.Args[0], err)
		return
	}

	// Create the interpreter
	interpreter := interpreter.New(string(data))

	// Evaluate the input
	out, err := interpreter.Evaluate(false)
	if err != nil {
		fmt.Printf("Error running program:%s\n", err)
		return
	}

	// Show the result
	fmt.Printf("\nResult:%s\n", out)
}
