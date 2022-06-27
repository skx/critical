// Set a variable
set a 43.1
puts "Variable a, ($$a), is set to: $a"

// variable expansion comes before execution.
set a pu
set b ts
$a$b "Hello World"

// expansion, once again.
// replacing things between the brackets with the output from executing them
puts [set a 4]
puts [set a]

// Variables can be longer.
set name "Steve Kemp"
puts "Hello World my name is $name"

// We have a standard library, located in `stdlib/stdlib.tcl`
//
// The standard library contains a couple of helpful methods,
// one of which is `assert_equal`.
//
// This will do "string" or "number" comparisons, and terminate
// execution on failure.
//
assert_equal "$name" "Steve Kemp"
assert_equal 9 [expr 3 * 3]
assert_equal 7.4 [- [+ 7 1.4] 1]
assert_equal 12 [+ 10 2]

// conditional
if { 1 } { puts "OK: 1 was ok" }
if { 0 } { puts "FAILURE: 0 was regarded as true" }
if { "steve" } { puts "OK: steve was ok" } else { puts "steve was not ok" }

// More conditionals
// remember we set some variables earlier:
//
//  "a" => "pu"
//  "b" => "ts"
//  "x" => UNDEFINED
//
if { $a } { puts "$$a is set" } else { puts "$$a is NOT set" }
if { $x } { puts "$$x is set - This is a bug" } else { puts "$$x is NOT set" }

//
// Setup some variables for a loop.
//
set i   1
set max 10
set sum 0

//
// Now we'll run a while-loop to sum some numbers
//
while { expr $i <= $max } {
   puts "  Loop $i"
   incr sum $i
   incr i
}

// Show the sum
puts "Sum of 1..10 (==(10x11)/2): $sum"




//
// Our first user-defined function!
//
proc inc {x} { puts "$$x is $x"; expr $x + 1 }
puts "3 inc is [inc 3]"

//
// Naive/Recursive solution.
//
proc fib {x} {
    if { expr $x <= 1 } {
        return 1
    } else {
        return [expr [fib [expr $x - 1]] + [fib [expr $x - 2]]]
    }
}

//
// A better, non-recursive, solution.
//
proc fib2 {n} {
    set a 1
    set b 1
    for {set N $n} {expr [set N] > 0} {decr N} {
        set tmp [+ $a $b]
        set a $b
        set b $tmp
    }
    return $a
}



//
// Lets run this in a loop
//
set i 0
set max 15

while { expr $i <= $max } {
    puts "Fib from a while-loop, with recursion, result $i is [fib $i]"
    incr i
}

//
// We can do the same thing again, using a for-loop, and a different (faster)
// Fibonacci sequence generator.
//
for {set i 0} {< $i 50} {incr i} {
    puts "Fib from a for-loop, without recursion, result $i is [fib2 $i]"
}


//
// This is just a horrid approach for running eval
//
set a { set b 20 ; incr b ; incr b; puts "$$b is $b" }
eval "$a"

//
// Is this better?
//
eval { set b 20 ; incr b ; incr b; puts "$$b is $b" }
