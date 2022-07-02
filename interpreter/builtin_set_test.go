package interpreter

import "testing"

func TestSet(t *testing.T) {

	// Create an empty interpreter
	e, er := New("")
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

	out, err := e.Evaluate()

	if err != nil {
		t.Fatalf("unexpected error running")
	}
	if out != "" {
		t.Fatalf("unexpected output")
	}

	// get the contents of a variable which is missing
	out, err = set(e, []string{"NAME"})
	if err != nil {
		t.Fatalf("error getting $NAME")
	}
	if out != "" {
		t.Fatalf("unexpected variable value")
	}

	// Now set the variable
	// get the contents of a variable which is missing
	out, err = set(e, []string{"NAME", "VALUE"})
	if err != nil {
		t.Fatalf("error getting $NAME")
	}
	if out != "VALUE" {
		t.Fatalf("unexpected variable value")
	}
}
