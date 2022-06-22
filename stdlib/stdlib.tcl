//
// This is the "standard library".
//
// We define some functions here which are available to all users
// of our application/scripting language.
//


//
// Maths functions should be easier to use.
//
// So we can write:
//
//    while { <= $a 5 } { .. }
//
// instead of
//
//    while { expr $a <= 5 } { ... }
//

proc + {a b} {
    expr $a + $b
}
proc - {a b} {
    expr $a - $b
}
proc / {a b} {
    expr $a / $b
}
proc * {a b} {
    expr $a * $b
}

//
// Comparison functions
//
proc < {a b} {
    expr $a < $b
}
proc <= {a b} {
    expr $a <= $b
}
proc > {a b} {
    expr $a > $b
}
proc >= {a b} {
    expr $a >= $b
}

//
// Equality
//
proc == {a b} {
    expr $a == $b
}
proc eq {a b} {
    expr $a eq $b
}
proc ne {a b} {
    expr $a ne $b
}


// Assert a condition is true.
proc assert {a b c} {
    if { expr $a $b $c } {
        puts "OK : $a $b $c"
    } else {
        puts "ERR: $a $b $c"
        exit 1
    }
}

// Assert two strings/numbers are equal.
proc assert_equal {a b} {

    // Is the first argument a number?
    if { regexp {^([0-9\.]+)$$} $a } {

        // Is the second argument a number?
        if { regexp {^([0-9\.]+)$$} $b } {

            // both numbers: numeric comparison
            return [ assert $a == $b ]
        } else {

            // $a is number
            // $b is string
            // -> string compare
            return [ assert $a eq $b ]
        }
    } else {
        // $a is string
        // $b is unknown
        // -> string compare
        return [ assert $a eq $b ]
    }
}

//
// Utility functions
//
proc repeat {n body} {
    set res ""
    while {> $n 0} {
        decr n
        set res [$body]
    }
    $res
}


//
// Now that we have a "repeat" function defined we could use it like so:
//
//   repeat 5 {
//        puts "Hello I'm alive";
//   }
//
// This would also work:
//
//   set foo 12
//   repeat 5 { incr foo }
//   => foo is now 17 (i.e. 12 + 5)
//
