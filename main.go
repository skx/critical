package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skx/critical/interpreter"
	"github.com/skx/critical/stdlib"
)

var version = "unreleased"

func main() {

	noStdlib := flag.Bool("no-stdlib", false, "Disable the (embedded) standard library.")
	versionFlag := flag.Bool("version", false, "Show our version, and exit.")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("critical %s\n", version)
		return
	}

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
	var out string
	var i *interpreter.Interpreter

	i, err = interpreter.New(input)
	if err != nil {
		fmt.Printf("Error creating interpreter %s\n", err)
		return
	}

	// Evaluate the input
	out, err = i.Evaluate()

	if err != nil && err != interpreter.ErrReturn {
		fmt.Printf("Error running program:%s\n", err)
		return
	}

	// Show the result
	fmt.Printf("\nResult:%s\n", out)
}
