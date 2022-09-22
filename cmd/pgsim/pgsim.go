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

	analysis, res, voltage, current, capasters, inductance, nodes := dios.Extract(*inFilePath)

	prettier.Info(map[string]interface{}{
		"1. Input File: ":      *inFilePath,
		"2. Precicion: ":       *e,
		"3. Max steps: ":       *maxSteps,
		"4. Nodes: ":           len(nodes),
		"5. Resistors: ":       len(res),
		"6. Current sources: ": len(current),
		"7. Voltage sources: ": len(voltage),
		"8. Capasters: ":       len(capasters),
		"9. Inductance: ":      len(inductance),
	})

	nodeBasedModel := model.NewModel(voltage, current, inductance, capasters, nodes, analysis, *maxSteps, *e)

	nodeBasedModel.Modeling()

	dios.WriteLogs(nodes, res, voltage, *inFilePath, *resultsPath)

	prettier.End(timer)

}
