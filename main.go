package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skx/critical/interpreter"
	"github.com/skx/critical/stdlib"
)

func main() {

	noStdlib := flag.Bool("no-stdlib", false, "Disable the (embedded) standard library")
	flag.Parse()

	// Ensure we have a file to execute.
	if len(flag.Args()) < 1 {
		fmt.Printf("Usage: critical file.tcl\n")
		return
	}

	// Read our standard library
	stdlib := stdlib.Contents()

	// Read the file the user wanted
	data, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		fmt.Printf("error reading file %s:%s\n", os.Args[0], err)
		return
	}

	// Join the two inputs, unless we shouldn't.
	input := string(data)
	if !*noStdlib {
		input = string(stdlib) + "\n" + input
	}

	// Create the interpreter
	interpreter := interpreter.New(input)

	// Evaluate the input
	out, err := interpreter.Evaluate()
	if err != nil {
		fmt.Printf("Error running program:%s\n", err)
		return
	}

	// Show the result
	fmt.Printf("\nResult:%s\n", out)
}
