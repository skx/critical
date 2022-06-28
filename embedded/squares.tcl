// Draw a simple square
proc draw_square {x y length} {
    move $x $y
    forwards $length
    turn 90
    forwards $length
    turn 90
    forwards $length
    turn 90
    forwards $length
}


// Pen down, so we draw.  (Default behaviour)
pen 1


// Draw a series of squares, rotating around in a circle
set count 25

// Loop $count times
for { set i 0 } { expr [set i] <= [set count] } { incr i } {

    // Draw a square in the middle, of width/size 65
    draw_square 150 150 65

    // Turn so that we complete a circle, or an approximation of one
    turn [expr 360 / [set count]]
}

// Save the `.PNG` and `.GIF` files
save

// All done
exit 0
