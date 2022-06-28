// Implementation of words exported to FORTH.

package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/skx/critical/interpreter"
)

// set the direction, absolutely.
func direction(i *interpreter.Interpreter, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("direction requires a single (numeric) argument")
	}

	val, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return "", err
	}

	// degree -> radians
	val = val * (math.Pi / 180)
	g.Direction(val)

	return "", nil
}

// move forwards.
func forwards(i *interpreter.Interpreter, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("forwards requires a single (numeric) argument")
	}

	val, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return "", err
	}
	g.Forward(val)

	return "", nil
}

// pen teleports to x,y.
func move(i *interpreter.Interpreter, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("move requires a pair of (numeric) arguments")
	}

	x := 0.0
	y := 0.0
	var err error

	x, err = strconv.ParseFloat(args[0], 64)
	if err != nil {
		return "", err
	}
	y, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return "", err
	}
	g.Move(x, y)

	return "", nil
}

// Set the pen up/down
func pen(i *interpreter.Interpreter, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("pen requires a single (numeric) argument")
	}

	val, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	if val == 0 {
		g.PenUp()
	} else {
		g.PenDown()
	}
	return "", nil
}

// Save the image, and the animation.
func save(i *interpreter.Interpreter, args []string) (string, error) {

	// write the image (PNG)
	err := g.WriteImage("turtle.png")
	if err != nil {
		return "", err
	}

	// Write the animation (GIF).
	err = g.WriteAnimation("turtle.gif")
	if err != nil {
		return "", err
	}

	saved = true

	return "", nil
}

// turn
func turn(i *interpreter.Interpreter, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("turn requires a single (numeric) argument")
	}

	val, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return "", err
	}

	// degree -> radians
	val = val * (math.Pi / 180)
	g.Turn(val)

	return "", nil
}
