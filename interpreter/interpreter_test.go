package interpreter

import (
	"strings"
	"testing"
)

func TestExpandEval(t *testing.T) {

	x, er := New(`puts [ expr 1 + [ expr 2 + 3 ] ]`)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

	out, err := x.Evaluate()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	expected := "6"

	if out != expected {
		t.Fatalf("unexpected output '%s'!=%s", out, expected)
	}

	// Ensure that we don't lose characters
	x, er = New(`puts "[expr 3 + 3]ab"`)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

	out, err = x.Evaluate()

	if err != nil {
		t.Fatalf("unexpected error running script")
	}
	if out != "6ab" {
		t.Fatalf("output %s != 6ab", out)
	}
}

func TestExpandString(t *testing.T) {

	x, er := New(`set a pu ; set b ts ; $a$b "OK"`)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
	x, er := New(`
proc squared {x} { expr $x * $x }
proc cubed {x} { expr $x * [squared $x] }

puts [cubed 9]
`)

	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
	x, er = New(`
proc multiply {x y} { expr $x * $y }

multiply 2
`)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
	x, er = New(`
proc multiply {x y} { expr $x * $y }

multiply 2 "steve"
`)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
	x, er = New(`
proc star {x y} { return [expr $x * $y] }

star 2 19
`)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
	x, er := New(`[ expr 1 + 1 ]`)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
	e, er := New(simple)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
	e, er = New(ifTrue)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
	e, er = New(ifFalse)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
	e, er = New(nested)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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

	e, er = New(abort)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

	out, err = e.Evaluate()
	if err != ErrReturn {
		t.Fatalf("unexpected error:%s", err)
	}
	if out != "313" {
		t.Fatalf("unexpected return value: got %s", out)
	}
}

func TestUnknownWord(t *testing.T) {
	e, er := New("moi")
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

	_, err := e.Evaluate()
	if err == nil {
		t.Fatalf("expected an error, got none:%s", err)
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("got an error, wrong kind:%s", err)
	}
}

func TestRepeat(t *testing.T) {

	//
	// Source code of the script we're going to run
	//
	src := `puts 34`

	//
	// Now run the thing in a loop
	//
	for n := 0; n < 10; n++ {
		//
		// Create the execution.
		//
		i, er := New(src)
		if er != nil {
			t.Fatalf("unexpected error creating interpreter")
		}

		//
		// return values
		//
		var err error
		var out string

		out, err = i.Evaluate()
		if err != nil {
			t.Fatalf("got unexpected error:%s", err)
		}
		if out != "34" {
			t.Fatalf("unexpected output: '%s'", out)
		}
	}

}
