package main

import (
	"flag"
	"pgsolver/pkg/dios"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"time"
)

func main() {
	timer := time.Now()                                                 // Program timer
	inFilePath := flag.String("f", "../data/ibmpg1.spice", "Iput file") // Path to input file
	resultsPath := flag.String("o", "./results", "Results folder path")
	e := flag.Float64("p", 1e-8, "Precision of modeling")                // Accuracy of the Zeidele method
	maxSteps := flag.Int("ms", 100000, "Max count of steps in modeling") // Max ammount of step during modeling

	flag.Parse()

	prettier.Start("PGSim", "1.1.0", "Ilya Shafeev", "MIT")
	prettier.Info(map[string]interface{}{"Input File: ": *inFilePath, "Precicion: ": *e, "Max steps: ": *maxSteps})

	res, voltage, current, nodes := dios.Extract(*inFilePath)
	nodeBasedModel := model.NewModel(voltage, current, nodes, *maxSteps, *e)

	nodeBasedModel.Prepare()
	nodeBasedModel.Init()
	nodeBasedModel.Modeling()

	dios.WriteLogs(nodes, res, *inFilePath, *resultsPath)

	prettier.End(timer)

}
