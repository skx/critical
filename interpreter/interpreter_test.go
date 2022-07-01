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
		t.Fatalf("unexpected error:%s", err)
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

func TestReturnScope(t *testing.T) {

	// Simple return from a procedure
	simple := `
proc test {} {
    return 32
    puts "NOT REACHED"
    return 99
}

test
`
	e := New(simple)

	out, err := e.Evaluate()
	if err != nil {
		t.Fatalf("unexpected error:%s", err)
	}
	if out != "32" {
		t.Fatalf("unexpected return value: got %s", out)
	}

	// Return from an if
	ifTrue := `
proc test2 {} {
    if { 1 } { return 32 } else { return 93 }
}

test2
`
	e = New(ifTrue)

	out, err = e.Evaluate()
	if err != nil {
		t.Fatalf("unexpected error:%s", err)
	}
	if out != "32" {
		t.Fatalf("unexpected return value: got %s", out)
	}

	// Simple return from an if
	ifFalse := `
proc test3 {} {
    if { 0 } { return 32 } else { return 93 }
}

test3
`
	e = New(ifFalse)

	out, err = e.Evaluate()
	if err != nil {
		t.Fatalf("unexpected error:%s", err)
	}
	if out != "93" {
		t.Fatalf("unexpected return value: got %s", out)
	}

	// nested if
	nested := `
proc test4 {} {
    if { 1 } {
       if { 1 } {
           return 17
       }
    }
    return 200
}

test4
`
	e = New(nested)

	out, err = e.Evaluate()
	if err != nil {
		t.Fatalf("unexpected error:%s", err)
	}
	if out != "17" {
		t.Fatalf("unexpected return value: got %s", out)
	}

	// Don't abort
	abort := `
proc test5 {} {
   return "Steve"
}

test5
return 313
`

	e = New(abort)

	out, err = e.Evaluate()
	if err != ErrReturn {
		t.Fatalf("unexpected error:%s", err)
	}
	if out != "313" {
		t.Fatalf("unexpected return value: got %s", out)
	}
}

func TestUnknownWord(t *testing.T) {
	e := New("moi")

	_, err := e.Evaluate()
	if err == nil {
		t.Fatalf("expected an error, got none:%s", err)
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("got an error, wrong kind:%s", err)
	}
}
