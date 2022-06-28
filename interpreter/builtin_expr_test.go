package interpreter

import (
	"strings"
	"testing"
)

func TestExpr(t *testing.T) {

	type TestCase struct {
		Input  []string
		Output string
		Error  string
	}

	tests := []TestCase{
		// basic maths
		{Input: []string{"3", "+", "3"}, Output: "6"},
		{Input: []string{"3.1", "+", "3.3"}, Output: "6.400000"},
		{Input: []string{"3", "*", "3"}, Output: "9"},
		{Input: []string{"3.1", "*", "3"}, Output: "9.300000"},
		{Input: []string{"3", "/", "3"}, Output: "1"},
		{Input: []string{"3.1", "/", "3"}, Output: "1.033333"},
		{Input: []string{"3", "-", "2"}, Output: "1"},
		{Input: []string{"3.4", "-", "1.2"}, Output: "2.200000"},

		// >
		{Input: []string{"3", ">", "2"}, Output: "1"},
		{Input: []string{"3.3", ">", "2"}, Output: "1"},
		{Input: []string{"1", ">", "2"}, Output: "0"},
		{Input: []string{"0.5", ">", "2"}, Output: "0"},

		// >=
		{Input: []string{"3", ">=", "3"}, Output: "1"},
		{Input: []string{"3.1", ">=", "3"}, Output: "1"},
		{Input: []string{"1", ">=", "2"}, Output: "0"},
		{Input: []string{"0.74", ">=", "2"}, Output: "0"},

		// <
		{Input: []string{"3", "<", "3"}, Output: "0"},
		{Input: []string{"3.0", "<", "3.0"}, Output: "0"},
		{Input: []string{"1", "<", "2"}, Output: "1"},
		{Input: []string{"1.4", "<", "2"}, Output: "1"},

		// <=
		{Input: []string{"3", "<=", "3"}, Output: "1"},
		{Input: []string{"21", "<=", "2"}, Output: "0"},

		// ==
		{Input: []string{"3", "==", "3"}, Output: "1"},
		{Input: []string{"21", "==", "2"}, Output: "0"},

		// !=
		{Input: []string{"3", "!=", "3"}, Output: "0"},
		{Input: []string{"21", "!=", "2"}, Output: "1"},

		{Input: []string{"steve", "eq", "steve"}, Output: "1"},
		{Input: []string{"steve", "eq", "Steve"}, Output: "0"},

		{Input: []string{"Kemp", "ne", "kemp"}, Output: "1"},
		{Input: []string{"Kemp", "ne", "Kemp"}, Output: "0"},

		// %
		{Input: []string{"3", "%", "3"}, Output: "0"},
		{Input: []string{"8", "%", "3"}, Output: "2"},

		// errors
		{Input: []string{"steve", "+", "3"}, Output: "", Error: "strconv"},
		{Input: []string{"34", "+", "steve"}, Output: "", Error: "strconv"},
		{Input: []string{"33", "^", "11"}, Output: "", Error: "unknown operation"},
	}

	for _, test := range tests {

		e := New("")

		out, err := expr(e, test.Input)
		if err != nil {
			if test.Error == "" {
				t.Fatalf("error calling Eval(%s):%s", test.Input, err)
			} else {
				if !strings.Contains(err.Error(), test.Error) {
					t.Fatalf("expected error, got the wrong one %s", err)
				}
			}
		}
		if out != test.Output {
			t.Fatalf("unexpected output for Eval(%s) - got %s, but expected %s", test.Input, out, test.Output)
		}
	}
}
