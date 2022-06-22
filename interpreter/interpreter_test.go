package interpreter

import "testing"

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
