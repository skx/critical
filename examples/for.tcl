//
// This example shows using a `for` loop.
//
// We also demonstrate the use of `continue` to skip an iteration, and
// `break` to escape the loop entirely.
//


//
// Show numbers from 1-10, skipping 5.
//
for {set x 1} {<= $x 20} {incr x} {

    // skip five
    if {== $x 5 } {
        puts "\t** Skipped this iteration"
        continue
    }

    // The `for` loop would run from 1-20, but we
    // use `break` to escape the loop after ten.
    if {> $x 10} {
        break
    }

    // Show something.
    puts "x is $x"
}
