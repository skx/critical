package interpreter

import "testing"

func TestBreak(t *testing.T) {

	out, err := breakFn(nil, []string{})

	if out != "" {
		t.Fatalf("unexpected output")
	}
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err != errBreak {
		t.Fatalf("got an error, but the wrong one:%v", err)
	}
}
