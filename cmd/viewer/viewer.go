package main

import (
	"bufio"
	"flag"
	"image"
	"image/png"
	"os"
	"pgsolver/pkg/draw"
	"pgsolver/pkg/prettier"
	"strconv"
	"strings"
)

func opacity(str string, tar string) float64 {
	if str == tar {
		return 0.35
	}
	return 1
}

func main() {
	inFilePath := flag.String("i", "modeling.plot", "Path to plot file")

	flag.Parse()

	inFile, err := os.Open(*inFilePath)

	if err != nil {
		prettier.Error("File open error.", err)
	}

	scanner := bufio.NewScanner(inFile)

	var minVal float64
	var maxVal float64
	var size [2]int
	var canvas *image.RGBA

	for scanner.Scan() {
		line := scanner.Text()
		splitedLine := strings.Split(line, " ")

		// Get file infos
		if line[0] == '#' {
			if splitedLine[1] == "Size" {
				if entryX, err := strconv.ParseInt(splitedLine[2], 10, 0); err == nil {
					size[0] = int(entryX)
				} else {
					prettier.Error("Geting x of size error.", err)
				}

				if entryY, err := strconv.ParseInt(splitedLine[3], 10, 0); err == nil {
					size[1] = int(entryY)
				} else {
					prettier.Error("Geting y of size error.", err)
				}

				canvas = draw.NewCanvas(size)
			}
			if splitedLine[1] == "Max_Val" {
				if entryV, err := strconv.ParseFloat(splitedLine[2], 64); err == nil {
					maxVal = entryV
				} else {
					prettier.Error("Geting Max_Val error.", err)
				}
			}
			if splitedLine[1] == "Min_Val" {
				if entryV, err := strconv.ParseFloat(splitedLine[2], 64); err == nil {
					minVal = entryV
				} else {
					prettier.Error("Geting Min_Val error.", err)
				}
			}

		} else {
			if entrV, err := strconv.ParseFloat(splitedLine[4], 64); err == nil {
				var startPoint [2]int
				var endPoint [2]int

				startSplited := strings.Split(splitedLine[1], "_")

				if entryX, err := strconv.ParseInt(startSplited[1], 10, 0); err == nil {
					startPoint[0] = int(entryX)
				} else {
					prettier.Error("Geting X of resistors start error.", err)
				}

				if entryY, err := strconv.ParseInt(startSplited[2], 10, 0); err == nil {
					startPoint[1] = int(entryY)
				} else {
					prettier.Error("Geting Y of resistors start error.", err)
				}

				endSplited := strings.Split(splitedLine[2], "_")

				if entryX, err := strconv.ParseInt(endSplited[len(endSplited)-2], 10, 0); err == nil {
					endPoint[0] = int(entryX)
				} else {
					prettier.Error("Geting X of resistors end error.", err)
				}

				if entryY, err := strconv.ParseInt(endSplited[len(endSplited)-1], 10, 0); err == nil {
					endPoint[1] = int(entryY)
				} else {
					prettier.Error("Geting Y of resistors end error.", err)
				}

				alph := opacity(splitedLine[0], "GND")
				rectColor := draw.GetColor(maxVal, minVal, entrV, alph)

				draw.DrawRect(canvas, startPoint, endPoint, rectColor)
			}
		}

	}

	splitedPath := strings.Split(*inFilePath, "/")
	projName := strings.Split(splitedPath[len(splitedPath)-1], ".")[0]

	f, _ := os.Create("./results/" + projName + "/" + projName + ".png")
	png.Encode(f, canvas)
}
