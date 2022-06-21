package interpreter

import "testing"

func TestComment(t *testing.T) {

	out, err := comment(nil, []string{})

	if out != "" {
		t.Fatalf("unexpected output")
	}
	if err != nil {
		t.Fatalf("unexpected error")
	}
}
