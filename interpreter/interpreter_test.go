package interpreter

import (
	"strings"
	"testing"
)

func TestExpandEval(t *testing.T) {

	x := New(`puts [ expr 1 + [ expr 2 + 3 ] ]`)

	out, err := x.Evaluate()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	expected := "6"

	if out != expected {
		t.Fatalf("unexpected output '%s'!=%s", out, expected)
	}

	// Ensure that we don't lose characters
	x = New(`puts "[expr 3 + 3]ab"`)
	out, err = x.Evaluate()

	if err != nil {
		t.Fatalf("unexpected error running script")
	}
	if out != "6ab" {
		t.Fatalf("output %s != 6ab", out)
	}
}

func TestExpandString(t *testing.T) {

	x := New(`set a pu ; set b ts ; $a$b "OK"`)

	out, err := x.Evaluate()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}
	if out != "OK" {
		t.Fatalf("unexpected output")
	}

	// now we have "$a -> pu"
	// now we have "$b -> ts"
	out = x.expandString("A$$A$a$b$c CC")
	if out != "A$Aputs CC" {
		t.Fatalf("unexpected output expanding string '%s'", out)
	}
}

// Define a function, and call it
func TestUserFunction(t *testing.T) {

	// Define a function call it
	x := New(`
proc squared {x} { expr $x * $x }
proc cubed {x} { expr $x * [squared $x] }

puts [cubed 9]
`)

	out, err := x.Evaluate()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	expected := "729"

	if out != expected {
		t.Fatalf("unexpected output '%s'!=%s", out, expected)
	}

	//
	// Now call a function with an error - wrong argument count
	//
	x = New(`
proc multiply {x y} { expr $x * $y }

multiply 2
`)

	_, err = x.Evaluate()
	if err == nil {
		t.Fatalf("expected an error, but got none")
	}
	if !strings.Contains(err.Error(), "function argument mismatch") {
		t.Fatalf("got an error, but the wrong one %s", err)
	}

	//
	// Another function with an error, wrong type
	//
	x = New(`
proc multiply {x y} { expr $x * $y }

multiply 2 "steve"
`)

	_, err = x.Evaluate()
	if err == nil {
		t.Fatalf("expected an error, but got none")
	}
	if !strings.Contains(err.Error(), "strconv") {
		t.Fatalf("got an error, but the wrong one %s", err)
	}

	//
	// A function with an explicit return :)
	//
	x = New(`
proc star {x y} { return [expr $x * $y] }

star 2 19
`)

	out, err = x.Evaluate()
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if out != "38" {
		t.Fatalf("wrong result for multiplication")
	}
}

func TestInvalidType(t *testing.T) {

	// This makes no sense, because we're replacing something
	// which could just be expressed directly.
	//
	//  i.e. "FOO" is the same as "[ FOO ]"
	//
	x := New(`[ expr 1 + 1 ]`)

	_, err := x.Evaluate()
	if err == nil {
		t.Fatalf("expected error on illegal token, got none")
	}
}
