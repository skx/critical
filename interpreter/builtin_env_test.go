package interpreter

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {

	// Read $foo - which should be empty
	out, err := env(nil, []string{"foo"})
	if err != nil {
		t.Fatalf("unexpected error reading  env-variable foo")
	}
	if out != "" {
		t.Fatalf("expected empty env-variable foo - got '%s'", out)
	}

	// Set the variable
	os.Setenv("foo", "bar")

	// Now we should have contents
	out, err = env(nil, []string{"foo"})
	if err != nil {
		t.Fatalf("unexpected error reading  env-variable foo")
	}
	if out != "bar" {
		t.Fatalf("env-variable foo - got '%s' not 'bar'", out)
	}
}
