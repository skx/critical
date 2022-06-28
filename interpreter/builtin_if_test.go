package interpreter

import (
	"strings"
	"testing"
)

func TestIf(t *testing.T) {

	type TestCase struct {
		In  string
		Err string
		Out string
	}

	tests := []TestCase{
		{
			In:  `if { 1 } { "steve" }`,
			Out: "steve",
		},
		{
			In:  `if { 0 } { "steve" } else { "bob" }`,
			Out: "bob",
		},
		{
			In:  `if { 0 } { "steve" } `,
			Out: "",
		},
		{
			In:  `if { expr 3 + } { "steve" } `,
			Out: "",
			Err: " expr requires three arguments",
		},
	}

	for _, test := range tests {

		// Run the example
		e := New(test.In)
		out, err := e.Evaluate()

		// Is this error expected?
		if err != nil {
			if test.Err == "" {
				t.Fatalf("unexpected error running %s:%s", test.In, err)
				continue
			}

			if !strings.Contains(err.Error(), test.Err) {
				t.Fatalf("expected error, but didn't get a match")
			}
		}
		if out != test.Out {
			t.Fatalf("input(%s) gave %s not %s", test.In, out, test.Out)
		}
	}
}
