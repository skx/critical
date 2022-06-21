package interpreter

import "testing"

func TestContinue(t *testing.T) {

	out, err := continueFn(nil, []string{})

	if out != "" {
		t.Fatalf("unexpected output")
	}
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err != errContinue {
		t.Fatalf("got an error, but the wrong one:%v", err)
	}
}
