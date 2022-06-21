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
