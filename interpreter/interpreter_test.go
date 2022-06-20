package interpreter

import "testing"

func TestBasic(t *testing.T) {

	// Basic program that does nothing useful
	x := New(`expr 3 + 6`)

	out, err := x.Evaluate(false)
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	expected := "9"

	if out != expected {
		t.Fatalf("unexpected output '%s'!=%s", out, expected)
	}
}

func TestEval(t *testing.T) {

	type TestCase struct {
		Input  string
		Output string
	}

	tests := []TestCase{
		TestCase{Input: `"Steve"`, Output: `Steve`},
		TestCase{Input: `33`, Output: `33`},
		TestCase{Input: `expr 3 + 3`, Output: `6`},
		TestCase{Input: `expr 3 * 3`, Output: `9`},
		TestCase{Input: `expr 3 - 1`, Output: `2`},
		TestCase{Input: `expr 3 / 3`, Output: `1`},
		TestCase{Input: `expr 3 <= 3`, Output: `1`},
		TestCase{Input: `expr 4 <= 3`, Output: `0`},

		TestCase{Input: `expr 3 < 3`, Output: `0`},
		TestCase{Input: `expr 1 < 3`, Output: `1`},

		TestCase{Input: `expr 3 > 3`, Output: `0`},
		TestCase{Input: `expr 5 > 3`, Output: `1`},

		TestCase{Input: `expr 3 >= 3`, Output: `1`},
		TestCase{Input: `expr 1 >= 3`, Output: `0`},

		TestCase{Input: `expr 1 == 3`, Output: `0`},
		TestCase{Input: `expr 2 == 2`, Output: `1`},
	}

	for _, test := range tests {

		e := New("")

		out, err := e.Eval(test.Input)
		if err != nil {
			t.Fatalf("error calling Eval(%s):%s", test.Input, err)
		}
		if out != test.Output {
			t.Fatalf("unexpected output for Eval(%s) - got %s, but expected %s", test.Input, out, test.Output)
		}
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

func TestExpandEval(t *testing.T) {

	x := New(`puts [ expr 1 + [ expr 2 + 3 ] ]`)

	out, err := x.Evaluate(false)
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	expected := "6"

	if out != expected {
		t.Fatalf("unexpected output '%s'!=%s", out, expected)
	}
}
