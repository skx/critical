package interpreter

import (
	"strings"
	"testing"
)

func TestFor(t *testing.T) {

	// Basic loop
	basic := `
set sum 0

for {set i 0  } { expr $i < 10 } { incr i} {
   set sum [incr sum $i ]
}
$sum
`
	// Loop with a continue
	continueLoop := `
set sum 0

for {set i 0  } { expr $i < 100 } { incr i} {
   if { expr $i >= 10 } { continue }
   set sum [incr sum $i ]
}
puts $sum
`
	// Loop with a break
	breakLoop := `
set sum 0

for {set i 0  } { expr $i < 100 } { incr i} {
   if { expr $i >= 10 } { break }
   set sum [incr sum $i ]
}
puts $sum
`

	// Run the basic loop
	e, er := New(basic)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

	out, err := e.Evaluate()

	// Is this error expected?
	if err != nil {
		t.Fatalf("unexpected error running code:%s", err)
	}
	if out != "45" {
		t.Fatalf("basic-loop failed: got %s", out)
	}

	// Run the break loop
	e, er = New(breakLoop)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

	out, err = e.Evaluate()

	// Is this error expected?
	if err != nil {
		t.Fatalf("unexpected error running code:%s", err)
	}
	if out != "45" {
		t.Fatalf("break-loop failed: got '%s'", out)
	}

	// Run the continue loop
	e, er = New(continueLoop)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

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
		`for { " } { expr $i < 100 } { incr i} { puts $i }`,
		`for { set i 0 } { " } { incr i} { puts $i }`,
		`for { set i 0 } { expr $i < 100} {  " } { puts $i }`,
		`for { set i 0 } { expr $i < 100} { incr i } { " }`,
	}

	for _, test := range tests {
		e, er = New(test)
		if er != nil {
			t.Fatalf("unexpected error creating interpreter")
		}

		_, err = e.Evaluate()

		if err == nil {
			t.Fatalf("expected error, got none")
		}
		if !strings.Contains(err.Error(), "unterminated string") {
			t.Fatalf("got error, but wrong one:%s", err)
		}
	}
}
