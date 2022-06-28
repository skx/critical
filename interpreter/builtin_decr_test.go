package interpreter

import "testing"

func TestDecr(t *testing.T) {

	// Create an empty interpreter
	e := New("")
	out, err := e.Evaluate()

	if err != nil {
		t.Fatalf("unexpected error running")
	}
	if out != "" {
		t.Fatalf("unexpected output")
	}

	// Set the variable "steve" to 30
	out, err = set(e, []string{"steve", "30"})
	if out != "30" {
		t.Fatalf("set had the wrong result")
	}
	if err != nil {
		t.Fatalf("unexpected error setting steve->3")
	}

	// Now decrease it by one.
	out, err = decr(e, []string{"steve"})
	if out != "29" {
		t.Fatalf("decr had the wrong result")
	}
	if err != nil {
		t.Fatalf("unexpected error setting steve->29")
	}

	// Now decrease it by two
	out, err = decr(e, []string{"steve", "2"})
	if out != "27" {
		t.Fatalf("decr had the wrong result")
	}
	if err != nil {
		t.Fatalf("unexpected error setting steve->27")
	}

	// finally test decreasing by a random non-number
	out, err = decr(e, []string{"steve", "not-a-number"})
	if out != "" {
		t.Fatalf("decr had the wrong result")
	}
	if err == nil {
		t.Fatalf("should have had an error, got none")
	}

	// Decrease a variable that doesn't exist
	out, err = decr(e, []string{"kemp"})
	if out != "-1" {
		t.Fatalf("decr had the wrong result")
	}
	if err != nil {
		t.Fatalf("unexpected error decreasing")
	}

	// floating-point decrement
	out, err = set(e, []string{"steve", "3.1"})
	if out != "3.1" {
		t.Fatalf("set had the wrong result %s != 3.1", out)
	}
	if err != nil {
		t.Fatalf("unexpected error")
	}

	// Now decrease it by one.
	out, err = decr(e, []string{"steve"})
	if out != "2.100000" {
		t.Fatalf("decr had the wrong result: %s", out)
	}
	if err != nil {
		t.Fatalf("unexpected error")
	}
}
