package main

import (
	"bufio"
	"flag"
	"pgsolver/pkg/utils"
	"strings"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

func NewSprite(scene *core.Node, start [2]int, end [2]int, size [2]int, cl [3]float32) {
	width := 5
	height := 5

	if start[0] != end[0] {
		width = end[0] - start[0]
	}
	if start[1] != end[1] {
		height = end[1] - start[1]
	}

	position := [2]float32{float32((start[0] + end[0]) / 2), float32((start[1] + end[1]) / 2)}

	mat := material.NewStandard(&math32.Color{R: cl[0], G: cl[1], B: cl[2]})
	sprite := graphic.NewSprite(float32(width)*0.1, float32(height)*0.1, mat)

	sprite.SetPosition((position[0]-float32(size[0]/2))*0.1, (-position[1]+float32(size[1]/2))*0.1, 0)

	scene.Add(sprite)
}

func GetColor(max float64, min float64, val float64, alph float64) [3]float32 {
	h := max / 4

	if val < h {
		return [3]float32{0, float32(val / h * 1), float32(1)}
	} else {
		if val < h*2 {
			return [3]float32{0, float32(1), float32((1 - (val-h)/h) * 1)}
		} else {
			if val < h*3 {
				return [3]float32{float32((val - h*2) / h * 1), float32(1), 0}
			} else {
				return [3]float32{float32(1), float32((1 - (val - h*3)) / h * 1), 0}
			}
		}
	}
}

func main() {

	inFilePath := flag.String("i", "modeling.plot", "Path to plot file")

	flag.Parse()

	scanner := bufio.NewScanner(utils.OpenFile(*inFilePath))

	var minVal float64
	var maxVal float64
	var size [2]int

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 10)
	cam.SetFar(5000)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetFramebufferSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0, 0, 0, 1.0)

	for scanner.Scan() {
		splitedLine := strings.Split(scanner.Text(), " ")

		// Get file info
		if splitedLine[0] == "#" {
			switch splitedLine[1] {
			case "Size":
				size[0] = utils.ParseInt(splitedLine[2])
				size[1] = utils.ParseInt(splitedLine[3])
			case "Max_Val":
				maxVal = utils.ParseFloat(splitedLine[2])
			case "Min_Val":
				minVal = utils.ParseFloat(splitedLine[2])
			}
		} else {
			startSplited := strings.Split(splitedLine[1], "_")
			endSplited := strings.Split(splitedLine[2], "_")

			startPoint := [2]int{utils.ParseInt(startSplited[1]), utils.ParseInt(startSplited[2])}
			endPoint := [2]int{utils.ParseInt(endSplited[len(endSplited)-2]), utils.ParseInt(endSplited[len(endSplited)-1])}
			voltageDrop := utils.ParseFloat(splitedLine[4])
			rectColor := GetColor(maxVal, minVal, voltageDrop, utils.Opacity(splitedLine[0], "GND"))

			NewSprite(scene, startPoint, endPoint, size, rectColor)
		}
	}

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
