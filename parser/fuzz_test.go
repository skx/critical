//go:build go1.18
// +build go1.18

package parser

import (
	"strings"
	"testing"
)

func FuzzParser(f *testing.F) {

	// empty string
	f.Add([]byte(""))

	// simple entries
	f.Add([]byte("puts \"ok\""))
	f.Add([]byte("puts [expr 3 + 3]"))

	// Assignments
	f.Add([]byte(`set a 3`))
	f.Add([]byte(`set a 3.32`))
	f.Add([]byte(`set a 10-10`))
	f.Add([]byte(`set a`))

	// and with usage
	f.Add([]byte(`set a pu
set b ts
$a$b "Hello"`))

	// Known errors are listed here.
	//
	// The purpose of fuzzing is to find panics, or unexpected errors.
	//
	// Some programs are obviously invalid though, so we don't want to
	// report those known-bad things.
	known := []string{
		"Closing ']' without opening one",
		"Closing '}' without opening one",
		"strconv.ParseInt",
		"unterminated pair",
		"unterminated string",
		"'-' may only occur at the start of the number",
	}

	f.Fuzz(func(t *testing.T, input []byte) {

		p := New(string(input))
		_, err := p.Parse()
		if err != nil {

			// not found this as a false-positive
			found := false

			// does it look familiar?
			for _, v := range known {
				if strings.Contains(err.Error(), v) {
					found = true
				}
			}

			// raise an error
			if !found {
				t.Fatalf("error parsing %s:%s", input, err)
			}
		}
	})
}
