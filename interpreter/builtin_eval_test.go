package interpreter

import "testing"

func TestEval(t *testing.T) {

	type TestCase struct {
		Input  string
		Output string
	}

	tests := []TestCase{
		{Input: `"Steve"`, Output: `Steve`},
		{Input: `"33"`, Output: `33`},
		{Input: `expr 3 + 3`, Output: `6`},
		{Input: `expr 3 * 3`, Output: `9`},
		{Input: `expr 3 - 1`, Output: `2`},
		{Input: `expr 3 / 3`, Output: `1`},
		{Input: `expr 3 <= 3`, Output: `1`},
		{Input: `expr 4 <= 3`, Output: `0`},

		{Input: `expr 3 < 3`, Output: `0`},
		{Input: `expr 1 < 3`, Output: `1`},

		{Input: `expr 3 > 3`, Output: `0`},
		{Input: `expr 5 > 3`, Output: `1`},

		{Input: `expr 3 >= 3`, Output: `1`},
		{Input: `expr 1 >= 3`, Output: `0`},

		{Input: `expr 1 == 3`, Output: `0`},
		{Input: `expr 2 == 2`, Output: `1`},
	}

	for _, test := range tests {

		e, er := New("")
		if er != nil {
			t.Fatalf("unexpected error creating interpreter")
		}

		out, err := evalFn(e, []string{test.Input})
		if err != nil {
			t.Fatalf("error calling eval(%s):%s", test.Input, err)
		}
		if out != test.Output {
			t.Fatalf("unexpected output for Eval(%s) - got %s, but expected %s", test.Input, out, test.Output)
		}
	}

}
