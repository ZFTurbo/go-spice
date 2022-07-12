package main

import (
	"flag"
	"pgsolver/pkg/dios"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
)

func main() {
	consPrettier := prettier.NewPrettier()
	inFilePath := flag.String("f", "../data/ibmpg1.spice", "Iput file")      // Path to input file
	outFilePath := flag.String("o", "../out/ibmpg1.solution", "Output file") // Path to output file
	e := flag.Float64("p", 1e-8, "Precision of modeling")                    // Accuracy of the Zeidele method
	maxSteps := flag.Int("ms", 100000, "Max count of steps in modeling")     // Max ammount of step during modeling

	flag.Parse()

	consPrettier.Start("PGSim", "1.0.0", "Ilya Shafeev")
	consPrettier.Info(map[string]interface{}{"Input File: ": *inFilePath, "Output File: ": *outFilePath, "Precicion: ": *e, "Max steps: ": *maxSteps})
	consPrettier.SetTimer()

	voltage, current, nodes, err := dios.Extract(*inFilePath)

	if err != nil {
		consPrettier.Error("File extraction error.", err)
	} else {
		nodeBasedModel := model.NewModel(voltage, current, nodes, *maxSteps, *e)

		nodeBasedModel.Prepare()
		nodeBasedModel.Init()
		nodeBasedModel.Modeling()

		dios.WriteLogs(nodes, *outFilePath)

		consPrettier.End()
	}
}
