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
// Run that in a loop to show some examples
//
loop cur 1 10 { puts "\t$cur! -> [fact $cur]" }

//
// Run a body multiple times, using an index "cur"
//

loop idx 0 10 { puts "Index:$idx Min:$min Max:$max"  }
