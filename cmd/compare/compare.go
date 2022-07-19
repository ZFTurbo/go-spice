package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"pgsolver/pkg/prettier"
	"pgsolver/pkg/utils"
	"strings"
)

func main() {
	prettier.Start("Comparing files", "", "", "")

	modelingFilePath := flag.String("m", "modeling.solution", "Modeling file") // Path to modeling file
	solutionFilePath := flag.String("s", "solution.solution", "Solution file") // Path to solution file

	modelingNodes := make(map[string]float64) //Modeling nodes voltage value
	sumPercentage := 0.0                      // Sum of percentage differences
	maxDifference := 0.0

	flag.Parse()

	// Open modeling and solution file
	modelingFile := utils.OpenFile(*modelingFilePath)
	solutionFile := utils.OpenFile(*solutionFilePath)

	scannerModeling := bufio.NewScanner(modelingFile) // The modeling file scanner
	scannerSolution := bufio.NewScanner(solutionFile) // The solution file scanner

	// Collect nodes and its values from modeling file
	for scannerModeling.Scan() {
		splitedLine := strings.Split(scannerModeling.Text(), " ")
		modelingNodes[splitedLine[0]] = utils.ParseFloat(splitedLine[1])
	}

	// Stack difference of solution and modeling nodes
	for scannerSolution.Scan() {
			splitedLine := strings.Split(scannerSolution.Text(), " ")
			voltageVal := utils.ParseFloat(splitedLine[2])

			if entryNode, found := modelingNodes[splitedLine[0]]; found {
				if entryNode != 0 && voltageVal != 0 {
					if entryNode >= voltageVal {
						sumPercentage += math.Abs(entryNode-voltageVal) / entryNode
						if math.Abs(entryNode-voltageVal)/entryNode > maxDifference {
							maxDifference = math.Abs(entryNode-voltageVal) / entryNode
						}
					} else {
						sumPercentage += math.Abs(entryNode-voltageVal) / voltageVal
						if math.Abs(entryNode-voltageVal)/voltageVal > maxDifference {
							maxDifference = math.Abs(entryNode-voltageVal) / voltageVal
						}
					}
				}
			} else {
				fmt.Printf("\n%sMissing node: %s, value: %f%s\n", prettier.Red, splitedLine[0], voltageVal, prettier.Reset)
			}
	}

	modelingFile.Close()
	solutionFile.Close()

	fmt.Println()
	prettier.Info(map[string]interface{}{"Avg difference %: ": (sumPercentage / float64(len(modelingNodes))) * 100, "Max difference %: ": maxDifference * 100})
	prettier.End()
}
