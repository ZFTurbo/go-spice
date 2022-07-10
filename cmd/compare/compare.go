package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("\nComparing files...")

	modelingFilePath := flag.String("m", "modeling.solution", "Modeling file") // Path to modeling file
	solutionFilePath := flag.String("s", "solution.solution", "Solution file") // Path to solution file

	modelingNodes := make(map[string]float64) //Modeling nodes voltage value
	sumPercentage := 0.0                      // Sum of percentage differences

	flag.Parse()

	// Open modeling and solution file
	modelingFile, errM := os.Open(*modelingFilePath)
	solutionFile, errS := os.Open(*solutionFilePath)

	if errM != nil {
		log.Fatal("Error in modeling file.\n", errM)
	}
	if errS != nil {
		log.Fatal("Error in solution file\n", errS)
	}

	scannerModeling := bufio.NewScanner(modelingFile) // The modeling file scanner
	scannerSolution := bufio.NewScanner(solutionFile) // The solution file scanner

	// Collect nodes and its values from modeling file
	for scannerModeling.Scan() {
		splitedLine := strings.Split(scannerModeling.Text(), " ")
		if entryV, err := strconv.ParseFloat(splitedLine[1], 64); err == nil {
			modelingNodes[splitedLine[0]] = entryV
		} else {
			log.Fatal("Value error in modeling scanner.\n", err)
		}
	}

	// Stack difference of solution and modeling nodes
	for scannerSolution.Scan() {
		splitedLine := strings.Split(scannerSolution.Text(), " ")
		if entryV, err := strconv.ParseFloat(splitedLine[2], 64); err == nil {
			if entryNode, found := modelingNodes[splitedLine[0]]; found {
				if entryNode != 0 && entryV != 0 {
					if entryNode >= entryV {
						sumPercentage += math.Abs(entryNode-entryV) / entryNode
					} else {
						sumPercentage += math.Abs(entryNode-entryV) / entryV
					}
				}
			} else {
				fmt.Printf("Missing node: %s, value: %f", splitedLine[0], entryV)
			}
		} else {
			log.Fatal("Value error in solution scanner.\n", err)
		}
	}

	difference := sumPercentage / float64(len(modelingNodes)) * 100 // Calculate difference of two files

	fmt.Printf("\nPercentage difference: %f%%\n\n", difference)

	modelingFile.Close()
	solutionFile.Close()
}
