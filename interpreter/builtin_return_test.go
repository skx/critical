package interpreter

import "testing"

func TestReturn(t *testing.T) {

	out, err := returnFn(nil, []string{"one"})

	if out != "one" {
		t.Fatalf("unexpected output")
	}
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err != errReturn {
		t.Fatalf("got an error, but the wrong one:%v", err)
	}
}
