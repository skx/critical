package interpreter

import (
	"testing"
)

// Benchmark_simple_return - This benchmark shows the overhead of just
// invoking a single `return` word.
//
func Benchmark_simple_return(b *testing.B) {

	//
	// Source code of the script we're going to run
	//
	src := `return 34`

	//
	// Create the interpreter
	//
	i, er := New(src)
	if er != nil {
		b.Fatalf("unexpected error creating interpreter")
	}

	//
	// return values
	//
	var err error
	var out string

	//
	// Now run the thing in a loop
	//
	b.ResetTimer()
	for n := 0; n < b.N; n++ {

		// execute
		out, err = i.Evaluate()

		// Ensure that the results are what we expect
		//
		// NOTE: This will kill the benchmark throughput a little.
		//
		if err != nil && err != ErrReturn {
			b.Fatal(err)
		}
		if out != "34" {
			b.Fail()
		}

	}
	b.StopTimer()

	//
	// Final execution is tested again.
	//
	if err != nil && err != ErrReturn {
		b.Fatal(err)
	}
	if out != "34" {
		b.Fail()
	}
}
