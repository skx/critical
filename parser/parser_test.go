package parser

import "testing"

func TestPuts(t *testing.T) {
	input := `puts "OK"`

	p := New(input)
	out, err := p.Parse()
	if err != nil {
		t.Fatalf("error parsing %s:%s", input, err)
	}

	if len(out) != 1 {
		t.Fatalf("wrong number of statements")
	}
}

func TestSeperator(t *testing.T) {
	input := `puts "Hello";

puts "World"`

	p := New(input)
	out, err := p.Parse()
	if err != nil {
		t.Fatalf("error parsing %s:%s", input, err)
	}

	if len(out) != 2 {
		t.Fatalf("wrong number of statements")
	}
}
