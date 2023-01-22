package dios

import (
	"strconv"
	"strings"

	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"pgsolver/pkg/utils"
)

// Writes solution in file
func WriteLogs(nodes map[string]*model.Node, filePath string, resultsPath string) {
	splitedPath := strings.Split(filePath, "/")
	projName := strings.Split(splitedPath[len(splitedPath)-1], ".")[0]

	utils.CreateFolder(resultsPath)
	utils.CreateFolder(resultsPath + "/" + projName)

	bar := prettier.DefaultBar(len(nodes), "Writing logs...")
	outFileNodeVal := utils.CreateFile(resultsPath + "/" + projName + "/" + projName + ".csv")

	outFileNodeVal.WriteString("Node, Value\n")

	for key, nodeInstance := range nodes {
		outFileNodeVal.WriteString(key + ", " + strconv.FormatFloat(nodeInstance.GetVoltage(), 'e', 8, 64) + "\n")
		bar.Add(1)
	}

	outFileNodeVal.Close()
}
