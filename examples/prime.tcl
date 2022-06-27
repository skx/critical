//
// Show prime numbers < 30
//
// The algorithm here is naive, but that's beside the point.  We're mostly
// testing loops and the modulus operation
//


proc is_prime {x} {

    // Less than two?  Not a prime
    if {< $x 2} {
        return "0"
    }

    // Search from N->X
    //
    // If the number is divisible by any of those values
    // then it cannot be prime.
    //
    // As an optimization we'd usually stop searching at SQRT(X)
    // but we don't have that primitive..
    //
    for { set n 2} { < [set n] [set x] } { incr n } {
        if { == [expr $x % $n] 0 } {
            return "0"
        }
    }
    return "1"
}

for {set i 0 } { <= [set i] 30 } { incr i } {
    if { is_prime $i } {
        puts "$i is prime"
    } else {
        puts "$i"
    }
}
