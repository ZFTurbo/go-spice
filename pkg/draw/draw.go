package draw

import (
	"image"
	"image/color"
)

// Create new color for heat map.
func GetColor(max float64, min float64, val float64, alph float64) color.NRGBA {
	h := max / 4

	if val < h {
		return color.NRGBA{0, uint8(val / h * 255 * alph), uint8(255 * alph), 255}
	} else {
		if val < h*2 {
			return color.NRGBA{0, uint8(255 * alph), uint8((1 - (val-h)/h) * 255 * alph), 255}
		} else {
			if val < h*3 {
				return color.NRGBA{uint8((val - h*2) / h * 255 * alph), uint8(255 * alph), 0, 255}
			} else {
				return color.NRGBA{uint8(255 * alph), uint8((1 - (val - h*3)) / h * 255 * alph), 0, 255}
			}
		}
	}
}

// Craete new canvas for image.
func NewCanvas(size [2]int) *image.RGBA {
	// Add padding to image
	size[0] += 10
	size[1] += 10

	canvas := image.NewRGBA(image.Rectangle{image.Point{}, image.Point{size[0], size[1]}})

	for x := 0; x < size[0]; x++ {
		for y := 0; y < size[1]; y++ {
			canvas.Set(x, y, color.Black)
		}
	}

	return canvas
}

// Draw new rectangle in canvas from start to end points filled by provided color.
func DrawRect(canvas *image.RGBA, start [2]int, end [2]int, cl color.Color) {
	widthCoords := [2]int{start[0] + 4, end[0] + 6}
	heightCoords := [2]int{start[1] + 4, end[1] + 6}

	for x := widthCoords[0]; x < widthCoords[1]+1; x++ {
		for y := heightCoords[0]; y < heightCoords[1]+1; y++ {
			canvas.Set(x, y, cl)
		}
	}
}
