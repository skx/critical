package interpreter

import (
	"testing"
)

func TestExit(t *testing.T) {

	// exit as a top-level
	exit := `exit 32`

	// exit within a loop
	basic := `
set sum 0
set i 0

while { expr $i < 10  } {
   set sum [incr sum $i ]
   exit 321
   incr i
}
puts $sum
`

	// exit within a proc
	proc := `
proc foo {a} {
   exit 43
}
foo 17
`

	// Run the exit-program
	e := New(exit)
	out, err := e.Evaluate()

	// Is this error expected?
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if err != ErrExit {
		t.Fatalf("unexpected error running code:%s", err)
	}
	if out != "32" {
		t.Fatalf("exit value didn't match - got '%s'", out)
	}

	// Run the loop-program
	e = New(basic)
	out, err = e.Evaluate()

	// Is this error expected?
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if err != ErrExit {
		t.Fatalf("unexpected error running code:%s", err)
	}
	if out != "321" {
		t.Fatalf("exit value didn't match - got '%s'", out)
	}

	// Run the proc-program
	e = New(proc)
	out, err = e.Evaluate()

	// Is this error expected?
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if err != ErrExit {
		t.Fatalf("unexpected error running code:%s", err)
	}
	if out != "43" {
		t.Fatalf("exit value didn't match - got '%s'", out)
	}
}
