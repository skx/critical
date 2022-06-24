//
// Test that we can return from functions, and that later
// expressions aren't executed
//
proc foo {a} {

    if {== 1 $a } {
        return 1
    } else {
        return 2
    }

    puts "This is a bug - this shouldn't be executed"
    return 3
}

//
// Test both the expected return-values
//
assert [foo 1] == 1
assert [foo 3] == 2
