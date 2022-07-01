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
proc % {a b} {
    expr $a % $b
}
proc ** {a b} {
    expr $a ** $b
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

//
// Min / Max
//
proc min {a b} {
    if {< $a $b } {
        return $a
    } else {
        return $b
    }
}

proc max {a b} {
    if {> $a $b} {
        return $a
    } else {
        return $b
    }
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
        }
    }

    // string compare
    assert $a eq $b
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



//
// Run a body multiple times, using an named variable for the index.
//
// You could use this like so:
//
//    loop cur 0 10 { puts "current iteration $cur ($min->$max)" }
//    => current iteration 0 (0-10)
//    => current iteration 1 (0-10)
//    ..
//    => current iteration 10 (0-10)
//
proc loop {var min max bdy} {
    // result
    set res ""

    // set the variable
    eval "set $var [set min]"

    // Run the test
    while {<= [set "$$var"] $max } {
        set res [$bdy]

        // This is a bit horrid
        eval {incr "$var"}
    }

    // return the last result
    $res
}
