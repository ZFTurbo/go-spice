package main

import (
	"flag"
	"time"

	"pgsolver/pkg/dios"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
)

func main() {
	timer := time.Now()                                                 // Program timer
	inFilePath := flag.String("f", "../data/ibmpg1.spice", "Iput file") // Path to input file
	resultsPath := flag.String("o", "./results", "Results folder path")
	e := flag.Float64("p", 1e-8, "Precision of modeling")                // Accuracy of the Zeidele method
	maxSteps := flag.Int("ms", 100000, "Max count of steps in modeling") // Max ammount of step during modeling

	flag.Parse()

	prettier.Start("PGSim", "1.1.0", "Ilya Shafeev", "MIT")

	nodes := dios.Extract(*inFilePath)

	prettier.Info(map[string]interface{}{
		"1. Input File: ":      *inFilePath,
		"2. Precicion: ":       *e,
		"3. Max steps: ":       *maxSteps,
		"4. Nodes: ":           len(nodes),
	})

	nodeBasedModel := model.NewModel(nodes, *maxSteps, *e)

	model.DCModeling(nodeBasedModel)

	dios.WriteLogs(nodes, *inFilePath, *resultsPath)

	prettier.End(timer)
}
