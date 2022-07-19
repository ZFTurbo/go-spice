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

	res, voltage, current, nodes := dios.Extract(*inFilePath)

	prettier.Info(map[string]interface{}{
		"1. Input File: ":         *inFilePath,
		"2. Precicion: ":          *e,
		"3. Max steps: ":          *maxSteps,
		"4. Nodes to be solved: ": len(nodes),
		"5. Resistors: ":          len(res),
		"6. Current sources: ":    len(current),
		"7. Voltage sources: ":    len(voltage),
	})

	nodeBasedModel := model.NewModel(voltage, current, nodes, *maxSteps, *e)

	nodeBasedModel.Prepare()
	nodeBasedModel.Init()
	nodeBasedModel.Modeling()

	dios.WriteLogs(nodes, res, voltage, *inFilePath, *resultsPath)

	prettier.End(timer)

}
