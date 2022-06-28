package interpreter

import "testing"

func TestIncr(t *testing.T) {

	// Create an empty interpreter
	e := New("")
	out, err := e.Evaluate()

	if err != nil {
		t.Fatalf("unexpected error running")
	}
	if out != "" {
		t.Fatalf("unexpected output")
	}

	// Set the variable "steve" to 3
	out, err = set(e, []string{"steve", "3"})
	if out != "3" {
		t.Fatalf("set had the wrong result")
	}
	if err != nil {
		t.Fatalf("unexpected error setting steve->3")
	}

	// Now increase it by one.
	out, err = incr(e, []string{"steve"})
	if out != "4" {
		t.Fatalf("incr had the wrong result")
	}
	if err != nil {
		t.Fatalf("unexpected error setting steve->4")
	}

	// Now increase it by two
	out, err = incr(e, []string{"steve", "2"})
	if out != "6" {
		t.Fatalf("incr had the wrong result")
	}
	if err != nil {
		t.Fatalf("unexpected error setting steve->6")
	}

	// finally test increasing by a random non-number
	// Now increase it by two
	out, err = incr(e, []string{"steve", "not-a-number"})
	if out != "" {
		t.Fatalf("incr had the wrong result")
	}
	if err == nil {
		t.Fatalf("should have had an error, got none")
	}

	// Increase a variable that doesn't exist
	out, err = incr(e, []string{"kemp"})
	if out != "1" {
		t.Fatalf("incr had the wrong result")
	}
	if err != nil {
		t.Fatalf("unexpected error increasing")

	}

	// floating-point increment
	out, err = set(e, []string{"steve", "3.1"})
	if out != "3.1" {
		t.Fatalf("set had the wrong result %s != 3.1", out)
	}
	if err != nil {
		t.Fatalf("unexpected error setting steve->3")
	}

	// Now increase it by one.
	out, err = incr(e, []string{"steve"})
	if out != "4.100000" {
		t.Fatalf("incr had the wrong result: %s", out)
	}
	if err != nil {
		t.Fatalf("unexpected error")
	}

}
