package interpreter

import "testing"

func TestRegexp(t *testing.T) {

	e, er := New("")
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

	// Test a valid match
	out, err := regexpFn(e, []string{
		"f.*d",
		"food",
	})
	if err != nil {
		t.Fatalf("unexpected error running regexp match:%s", err)
	}
	if out != "1" {
		t.Fatalf("unexpected regexp match result")
	}

	// Test a failed match
	out, err = regexpFn(e, []string{
		"f.*d",
		"cake",
	})
	if err != nil {
		t.Fatalf("unexpected error running regexp match:%s", err)
	}
	if out != "0" {
		t.Fatalf("unexpected regexp match result")
	}

	// Test a bogus regexp
	_, err = regexpFn(e, []string{
		"*",
		"cake",
	})
	if err == nil {
		t.Fatalf("expected an error compiling a bogus regexp, got none")
	}

}
