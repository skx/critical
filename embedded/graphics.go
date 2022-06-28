package main

import (
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math"
	"os"

	"github.com/andybons/gogif"
)

// Graphics holds our state.
type Graphics struct {
	Image       *image.RGBA
	Pos         Position
	Orientation float64
	Color       color.Color
	Pen         bool
	Gif         *gif.GIF
}

// Position holds the coordinates of our pen.
type Position struct {
	X, Y float64
}

// NewGraphics creates a new object.
func NewGraphics(i *image.RGBA, starting Position) (t *Graphics) {
	t = &Graphics{
		Image:       i,
		Pos:         starting,
		Orientation: 0.0,
		Color:       color.Black,
		Pen:         true,
		Gif:         &gif.GIF{},
	}

	return
}

// Move updates the location of the pen.
func (t *Graphics) Move(x float64, y float64) {
	t.Pos.X = x
	t.Pos.Y = y
}

// Forward moves the pen forwards, in a straight-line, from the
// current location - in the direction of travel.
func (t *Graphics) Forward(dist float64) {
	for i := 0; i < int(dist); i++ {
		if t.Pen {
			t.Image.Set(int(t.Pos.X), int(t.Pos.Y), t.Color)
		}

		x := 1.0 * math.Sin(t.Orientation)
		y := 1.0 * -math.Cos(t.Orientation)

		t.Pos = Position{t.Pos.X + x, t.Pos.Y + y}
	}

	//
	// Append to our animation.
	//
	// Need to convert from RGBA -> Paletted Image
	if t.Pen {
		bounds := t.Image.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: 64}
		quantizer.Quantize(palettedImage, bounds, t.Image, image.Point{})

		// Append the new frame
		t.Gif.Image = append(t.Gif.Image, palettedImage)
		t.Gif.Delay = append(t.Gif.Delay, 0)
	}
}

// Turn adds the specified number of degrees to our direction.
func (t *Graphics) Turn(radians float64) {
	t.Orientation += radians
}

// Direction sets the absolute direction of the pen.
func (t *Graphics) Direction(radians float64) {
	t.Orientation = radians
}

// PenUp lifts the pen, so movement will not draw anything.
func (t *Graphics) PenUp() {
	t.Pen = false
}

// PenDown lowers the pen, so movement will draw.
func (t *Graphics) PenDown() {
	t.Pen = true
}

// WriteImage outputs the end result - as a PNG
func (t *Graphics) WriteImage(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, t.Image)

	return err
}

// WriteAnimation outputs the end result - as a GIF
func (t *Graphics) WriteAnimation(name string) error {

	// write the gif
	h, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer h.Close()
	err = gif.EncodeAll(h, t.Gif)

	return err
}
