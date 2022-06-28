proc draw_square {x y length} {
    puts "Moving to $x $y"
    move $x $y
    pen 1
    forwards $length
    turn 90
    forwards $length
    turn 90
    forwards $length
    turn 90
    forwards $length
}

puts "Starting .."
draw_square 50 50 25
puts "Saving .."
save
puts "Done"
