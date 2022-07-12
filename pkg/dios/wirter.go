package dios

import (
	"os"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"strconv"
)

// Writes solution in file
func WriteLogs(nodes map[string]*model.Node, fileName string) {
	consPrettier := prettier.NewPrettier()
	bar := consPrettier.DefaultBar(len(nodes), "Writing logs...")

	outFile, err := os.Create(fileName)

	if err != nil {
		consPrettier.Error("Create output file error.", err)
	} else {
		for key, nodeInstance := range nodes {
			outFile.WriteString(key + " ")
			outFile.WriteString(strconv.FormatFloat(nodeInstance.V, 'e', 8, 64))
			outFile.WriteString("\n")
			bar.Add(1)
		}
	}

	bar.Close()
	outFile.Close()
}
