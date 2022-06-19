// variable
set a [ expr 2 + 2 ]
puts "A is set to: $a"

// variable expansion comes first.
set a pu
set b ts
$a$b "Hello World"

// expansion
puts [set a 4]
puts [set a]

// variable
set name Steve
puts "Hello World my name is $name"

// maths
puts [expr 3 * 4]  ok world
puts I'm dividing - NOTE this is an integer expression - [expr 3 / 4]

// conditional
if { 1 } { puts "OK: 1 was ok" }
if { 0 } { puts "FAILURE: 0 was regarded as true" }

if { "steve" } { puts "OK: steve was ok" } else { puts "steve was not ok" }

puts "I'm still alive, after the conditional execution!"

if { $a } { puts "A is set" } else { puts "A is NOT set" }
if { $x } { puts "X is set - This is a bug" } else { puts "X is NOT set" }

//
// Horrid
//
set i 1
set max 10

// Because we don't have a proper parser so we can't handle a newline inside
// our block, so we have to fake it via the use of the continuation-line
// which is triggered by a trailing "\" character.
//
// We need this because our blocks must be on "one line", but the puts
// command would consume the whole line.
//
// i.e. { puts "foo" incr bar } would actually invoke "puts" with the
// text '"foo" incr bar', meaning the increment would never execute and
// our loop would last forever.
while { expr $i <= $max } { \
   puts "  Loop $i" \
   incr i  \
}
