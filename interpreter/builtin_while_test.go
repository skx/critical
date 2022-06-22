package interpreter

import (
	"strings"
	"testing"
)

func TestWhile(t *testing.T) {

	// Basic loop
	basic := `
set sum 0
set i 0

while { expr $i < 10  } {
   set sum [incr sum $i ]
   incr i
}
puts $sum
`
	// Loop with a continue
	continueLoop := `
set sum 0
set i 0

while { expr $i < 100 }  {
   incr i
   if { expr $i >= 10 } { continue }
   set sum [incr sum $i ]
}
puts $sum
`
	// Loop with a break
	breakLoop := `
set sum 0
set i 0

while { expr $i < 100 } {
   if { expr $i >= 10 } { break }
   set sum [incr sum $i ]
   incr i
}
puts $sum
`

	// Run the basic loop
	e := New(basic)
	out, err := e.Evaluate()

	// Is this error expected?
	if err != nil {
		t.Fatalf("unexpected error running code:%s", err)
	}
	if out != "45" {
		t.Fatalf("basic-loop failed: got %s", out)
	}

	// Run the break loop
	e = New(breakLoop)
	out, err = e.Evaluate()

	// Is this error expected?
	if err != nil {
		t.Fatalf("unexpected error running code:%s", err)
	}
	if out != "45" {
		t.Fatalf("break-loop failed: got '%s'", out)
	}

	// Run the continue loop
	e = New(continueLoop)
	out, err = e.Evaluate()

	// Is this error expected?
	if err != nil {
		t.Fatalf("unexpected error running code:%s", err)
	}
	if out != "45" {
		t.Fatalf("continue-loop failed: got '%s'", out)
	}

	// Bogus loops
	tests := []string{
		// error in conditional
		`set i 0 ; while { " } { puts $i }`,

		// error in body
		`set i 0 ; while { expr $i < 100} { incr i ;  " }`,

		// error after the first loop.
		// i.e. retesting the conditional will fail the second
		// time as "while "steve" < 1000" will fail.
		`set i 0 ; while { expr $i < 100} { set i "steve" }`,
	}

	for _, test := range tests {
		e = New(test)
		_, err = e.Evaluate()

		if err == nil {
			t.Fatalf("expected error, got none")
		}
		if !strings.Contains(err.Error(), "unterminated string") && !strings.Contains(err.Error(), "strconv.") {
			t.Fatalf("got error, but wrong one:%s", err)
		}
	}
}
