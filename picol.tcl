proc square {x} {
    expr $x * $x
}

set a 1
while {expr $a <= 10} {
    if {expr $a == 5} {
        puts "\tMissing five!"
        incr a
        continue
        puts "After continue this won't be executed"
    }
    puts "I can compute that $a*$a = [square $a]"
    set a [expr $a + 1]
}

puts "I'm alive; after the 'while' loop."
