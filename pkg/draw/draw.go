package draw

import (
	"image"
)

func NewCanvas(size [2]int) *image.Rectangle {
	canvas := &image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: size[0], Y: size[1]}}
	return canvas
}


