//
// This example demonstrates the use of the `loop` word,
// which allows you to repeat a block with a start/end index
// using the named variable as index.
//



//
// Define a command to calculate a factorial, recursively.
//
proc fact {n} {
    if  {<= $n 1} {
        return 1
    } else {
        return [* $n [fact [- $n 1]]]
    }
}

//
// Run that in a loop to show some examples.
//
// Here the initial three parameters to the 'loop'
// word are the name of the variable to use, within
// the body, and the min/max indexes.
//
loop cur 1 10 { puts "\t$cur! -> [fact $cur]" }

//
// Run a body multiple times, using an index "idx"
//
// NOTE: Here we use "min" and "max" which are the
// the bounding values for the loop.
//
loop idx 0 10 { puts "Index:$idx Min:$min Max:$max"  }
