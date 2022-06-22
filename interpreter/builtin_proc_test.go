package interpreter

import "testing"

func TestProc(t *testing.T) {

	e := New("")

	_, ok := e.functions["squared"]
	if ok {
		t.Fatalf("function 'squared' already exists")
	}

	out, err := proc(e, []string{
		"squared",
		"{x}",
		"{expr $x * $x}"})

	if err != nil {
		t.Fatalf("error calling proc")
	}
	if out != "" {
		t.Fatalf("unexpected output from proc:%s", out)
	}

	// Now we should have the function
	_, ok = e.functions["squared"]
	if !ok {
		t.Fatalf("function 'squared' doesn't exist, after definition")
	}

}
