package interpreter

import "testing"

func TestBasic(t *testing.T) {

	// Basic program that does nothing useful
	x := New(`expr 3 + 6`)

	out, err := x.Evaluate(false)
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	expected := "9.000000"

	if out != expected {
		t.Fatalf("unexpected output '%s'!=%s", out, expected)
	}
}

func TestDecr(t *testing.T) {

	// Set a value, and then decrease it..
	x := New(`set a 10 ; decr a ; decr a 2; decr a; set a`)

	out, err := x.Evaluate(false)
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	expected := "6"

	if out != expected {
		t.Fatalf("unexpected output '%s'!=%s", out, expected)
	}
}

func TestIncr(t *testing.T) {

	// Basic program that increments an empty variable
	x := New(`incr a ; incr a 1; incr a`)

	out, err := x.Evaluate(false)
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	expected := "3"

	if out != expected {
		t.Fatalf("unexpected output '%s'!=%s", out, expected)
	}
}
