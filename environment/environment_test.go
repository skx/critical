package environment

import "testing"

// TestGetSet tests get/set on a variable
func TestGetSet(t *testing.T) {

	e := New()

	out := ""
	ok := false

	// by default the environment is empty
	_, ok = e.Get("FOO")
	if ok {
		t.Fatalf("fetching missing variable shouldn't work")
	}

	// Now set
	e.Set("FOO", "BAR")
	out, ok = e.Get("FOO")
	if !ok {
		t.Fatalf("fetching variable shouldn't fail")
	}
	if out != "BAR" {
		t.Fatalf("variable had wrong value")
	}

	// Clear
	e.Clear("FOO")

	// Fetching should now fail
	_, ok = e.Get("FOO")
	if ok {
		t.Fatalf("fetching missing variable shouldn't work")
	}
}

func TestScopedSet(t *testing.T) {

	// parent
	p := New()
	p.Set("FOO", "BAR")

	// child
	c := NewEnclosedEnvironment(p)

	// Child should be able to reach parent variable
	val, ok := c.Get("FOO")
	if !ok {
		t.Fatalf("failed to get variable in parent scope")
	}
	if val != "BAR" {
		t.Fatalf("got variable; wrong value")
	}

	// Set variable in child
	c.Set("NAME", "STEVE")

	// child should get it
	val, ok = c.Get("NAME")
	if !ok {
		t.Fatalf("failed to get child-variable in child")
	}
	if val != "STEVE" {
		t.Fatalf("variable had wrong value")
	}

	// parent should not
	_, ok = p.Get("NAME")
	if ok {
		t.Fatalf("shouldn't be able to get child-variable in parent")
	}

	// set in the child
	//
	// Will actually set in the parent
	c.Set("FOO", "BART")

	val, ok = p.Get("FOO")
	if !ok {
		t.Fatalf("parent-child weirdness")
	}
	if val != "BART" {
		t.Fatalf("parent-child set failed")
	}
}
