package interpreter

import "testing"

// TestArity tests calling built-in functions with the wrong number of args
func TestArity(t *testing.T) {

	tests := []string{
		`break "one"`,

		`continue "one" "two"`,

		`decr`,
		`decr "one" "two" "three"`,

		`exit`,
		`exit "one" "two"`,

		`expr 1`,
		`expr 1 + `,
		`expr 1 + 2 + 3`,

		`if { 1 } `,
		`if { 1 } { 2 } else { 3 } or { 4}`,

		`incr`,
		`incr "one" 2 3`,

		`proc "one"`,
		`proc "one", "two", "three", "four"`,

		`puts "One" "Two"`,
		`puts`,

		`return`,
		`return "one" "two"`,
		`set`,
		`set 1 2 3`,

		`while { 1 } `,
		`while { 1 } { 2 } { 3  }`,
	}

	for _, test := range tests {
		x := New(test)

		_, err := x.Evaluate()

		if err == nil {
			t.Fatalf("expected error, got none:%s", test)
		}
	}
}
