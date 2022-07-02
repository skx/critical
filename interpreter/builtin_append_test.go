package interpreter

import "testing"

func TestAppend(t *testing.T) {

	src := `
set var 0
for {set i 1} {expr $i <= 10} {incr i} {
    append var "," $i
}
$var
`
	// Run the example
	e, er := New(src)
	if er != nil {
		t.Fatalf("unexpected error creating interpreter")
	}

	out, err := e.Evaluate()

	// Is this error expected?
	if err != nil {
		t.Fatalf("unexpected error running code:%s", err)
	}
	if out != "0,1,2,3,4,5,6,7,8,9,10" {
		t.Fatalf("append-test failed: got %s", out)
	}
}
