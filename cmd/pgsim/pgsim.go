package main

import (
	"flag"
	"pgsolver/pkg/dios"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"time"
)

func main() {
	timer := time.Now() // Program timer
	inFilePath := flag.String("f", "../data/ibmpg1.spice", "Iput file")      // Path to input file
	outFilePath := flag.String("o", "../out/ibmpg1.solution", "Output file") // Path to output file
	e := flag.Float64("p", 1e-8, "Precision of modeling")                    // Accuracy of the Zeidele method
	maxSteps := flag.Int("ms", 100000, "Max count of steps in modeling")     // Max ammount of step during modeling

	flag.Parse()

	prettier.Start("PGSim", "1.0.5", "Ilya Shafeev")
	prettier.Info(map[string]interface{}{"Input File: ": *inFilePath, "Output File: ": *outFilePath, "Precicion: ": *e, "Max steps: ": *maxSteps})

	voltage, current, nodes, err := dios.Extract(*inFilePath)

	if err != nil {
		prettier.Error("File extraction error.", err)
	} else {
		nodeBasedModel := model.NewModel(voltage, current, nodes, *maxSteps, *e)

		nodeBasedModel.Prepare()
		nodeBasedModel.Init()
		nodeBasedModel.Modeling()

		dios.WriteLogs(nodes, *outFilePath)

		prettier.End(timer)
	}
}
