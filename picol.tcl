// Procedure (i.e. function) which squares the given number.
proc square {x} {
    expr $x * $x
}

set a 1
while {<= $a 10} {
    if {== $a 5} {
        puts "\tMissing five!"
        set a [+ $a 1]
        continue
        puts "After continue this won't be executed"
    }
    puts "I can compute that $a*$a = [square $a]"
    set a [+ $a 1]
}

puts "I'm alive; after the 'while' loop."
