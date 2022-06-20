// variable
set a [ expr 2 + 2 ]
puts "A is set to: $a"

// variable expansion comes before execution.
set a pu
set b ts
$a$b "Hello World"

// expansion, once again.
puts [set a 4]
puts [set a]

set name "Steve"
puts "Hello World my name is $name"

// maths
puts "[expr 3 * 4]  ok world"
puts "I'm dividing - NOTE this is an integer expression - [expr 3 / 4]"

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
if { $a } { puts "A is set" } else { puts "A is NOT set" }
if { $x } { puts "X is set - This is a bug" } else { puts "X is NOT set" }

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

puts "Sum of 1..10 (==(10x11)/2): $sum"


