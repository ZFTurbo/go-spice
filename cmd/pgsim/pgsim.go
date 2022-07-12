package main

import (
	"flag"
	"fmt"
	"os"
	"pgsolver/pkg/extractor"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"strconv"
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

	inExtactor := extractor.NewExtractor(*inFilePath)
	voltage, current, nodes, err := inExtactor.Extract()

	if err != nil {
		consPrettier.Error("File extraction error.", err)
	} else {
		nodeBasedModel := model.NewModel(voltage, current, nodes, *maxSteps, *e)

		nodeBasedModel.Prepare()
		nodeBasedModel.Init()
		nodeBasedModel.Modeling()

		fmt.Println("Writing results...")

		// Log out data
		outFile, err := os.Create(*outFilePath)
		if err != nil {
			consPrettier.Error("Create output file error.", err)
		} else {
			for key, nodeInstance := range nodes {
				outFile.WriteString(key + " ")
				outFile.WriteString(strconv.FormatFloat(nodeInstance.V, 'e', 8, 64))
				outFile.WriteString("\n")
			}
		}

		outFile.Close()

		consPrettier.End()
	}
}
