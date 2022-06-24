//
// This example demonstrates the use of the `append` and `for` words
//
// `append` appends strings to variables.
//
// `for` allows you to write loops.
//


// Start with zero
set var "0"

//
// Loop ten times adding numbers to the variable
//
for { set i 1 } { <= [set i] 10 } { incr i } {
    append var ",$i"
}

//
// Ensure we get the result we expect
//
assert_equal [set var] "0,1,2,3,4,5,6,7,8,9,10"
