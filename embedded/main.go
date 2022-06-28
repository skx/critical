// This is a simple application, allowing retro scripted graphics
// to be produced.
//
// We create a global image ("i") which we draw to, via the TCL words
// "forward", "move", etc.
//
// The end result is output as a PNG.  However we also generate a GIF of
// the drawing process - updating at each step of the drawing process.
//
// NOTE: This is very slow for circles, for obvious reasons..
//

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"

	"github.com/skx/critical/interpreter"
)

// Our embedded interpreter.
var e *interpreter.Interpreter

// Graphics-helper, for drawing into our image.
var g *Graphics

// Did an image get saved?
var saved bool

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: embedded path/to/file.tcl\n")
		return
	}

	// Read the file the user wanted
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("error reading file %s:%s\n", os.Args[0], err)
		return
	}

	// New PNG image - with white background
	i := image.NewRGBA(image.Rect(0, 0, 300, 300))
	c := color.RGBA{255, 255, 255, 255}
	draw.Draw(i, i.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)

	// Create a new helper to draw into our image.
	//
	// Point is set to the middle of the image.
	g = NewGraphics(i, Position{150.0, 150.0})

	// Create the interpreter instance
	fmt.Printf("Running script\n")
	e = interpreter.New(string(data))

	e.RegisterBuiltin("direction", direction)
	e.RegisterBuiltin("forwards", forwards)
	e.RegisterBuiltin("move", move)
	e.RegisterBuiltin("pen", pen)
	e.RegisterBuiltin("save", save)
	e.RegisterBuiltin("turn", turn)

	out, err := e.Evaluate()
	if err == interpreter.ErrExit || err == interpreter.ErrReturn {
		fmt.Printf("Script exited with either `return` or `exit`.\n")
		if saved {
			fmt.Printf("Image saved\n")
		} else {
			fmt.Printf("No image was saved - did you call 'save'?\n")
		}
		return
	}
	if err != nil {
		fmt.Printf("Error running script: %s\n", err)
	}

	if saved {
		fmt.Printf("Script exited cleanly with output '%s'\n", out)
	} else {
		fmt.Printf("Script exited cleanly, but no output was generated\n")
		fmt.Printf("Did you forget to call `save`?\n")
	}

}
