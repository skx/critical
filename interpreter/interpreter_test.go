package interpreter

import "testing"

func TestBasic(t *testing.T) {

	// Basic program that does nothing useful
	x := New(`expr 3 + 6`)

	out, err := x.Evaluate()
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
		{Input: `"Steve"`, Output: `Steve`},
		{Input: `"33"`, Output: `33`},
		{Input: `expr 3 + 3`, Output: `6`},
		{Input: `expr 3 * 3`, Output: `9`},
		{Input: `expr 3 - 1`, Output: `2`},
		{Input: `expr 3 / 3`, Output: `1`},
		{Input: `expr 3 <= 3`, Output: `1`},
		{Input: `expr 4 <= 3`, Output: `0`},

		{Input: `expr 3 < 3`, Output: `0`},
		{Input: `expr 1 < 3`, Output: `1`},

		{Input: `expr 3 > 3`, Output: `0`},
		{Input: `expr 5 > 3`, Output: `1`},

		{Input: `expr 3 >= 3`, Output: `1`},
		{Input: `expr 1 >= 3`, Output: `0`},

		{Input: `expr 1 == 3`, Output: `0`},
		{Input: `expr 2 == 2`, Output: `1`},
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
	out = x.expandString("AA$a$b$c CC")
	if out != "AAputs CC" {
		t.Fatalf("unexpected output expanding string '%s'", out)
	}
}
