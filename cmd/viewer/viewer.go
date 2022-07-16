package main

import (
	"bufio"
	"flag"
	"image"
	"image/png"
	"os"
	"pgsolver/pkg/draw"
	"pgsolver/pkg/utils"
	"strings"
)

func main() {
	inFilePath := flag.String("i", "modeling.plot", "Path to plot file")

	flag.Parse()

	scanner := bufio.NewScanner(utils.OpenFile(*inFilePath))

	var minVal float64
	var maxVal float64
	var size [2]int
	var canvas *image.RGBA

	for scanner.Scan() {
		splitedLine := strings.Split(scanner.Text(), " ")

		// Get file info
		if splitedLine[0] == "#" {
			switch splitedLine[1] {
			case "Size":
				size[0] = utils.ParseInt(splitedLine[2])
				size[1] = utils.ParseInt(splitedLine[3])
				canvas = draw.NewCanvas(size)
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
			rectColor := draw.GetColor(maxVal, minVal, voltageDrop, utils.Opacity(splitedLine[0], "GND"))

			draw.DrawRect(canvas, startPoint, endPoint, rectColor)
		}
	}

	splitedPath := strings.Split(*inFilePath, "/")
	projName := strings.Split(splitedPath[len(splitedPath)-1], ".")[0]
	resultsPath := strings.Split(*inFilePath, projName)[0]

	f, _ := os.Create(resultsPath + projName + "/" + projName + ".png")
	png.Encode(f, canvas)
}

func OpenFile(s string) {
	panic("unimplemented")
}
