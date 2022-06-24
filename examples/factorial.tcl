//
// Demonstrate using recursion to calculate factorials.
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
