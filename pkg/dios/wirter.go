package dios

import (
	"math"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"pgsolver/pkg/utils"
	"sort"
	"strconv"
	"strings"
)

// Writes solution in file
func WriteLogs(nodes map[string]*model.Node, res map[string]string, volage map[string]float64, filePath string, resultsPath string) {
	splitedPath := strings.Split(filePath, "/")
	projName := strings.Split(splitedPath[len(splitedPath)-1], ".")[0]

	utils.CreateFolder(resultsPath)
	utils.CreateFolder(resultsPath + "/" + projName)

	bar := prettier.DefaultBar(len(nodes)+2*len(res)+len(volage), "Writing logs...")
	outFileNodeVal := utils.CreateFile(resultsPath + "/" + projName + "/" + projName + ".solution")

	var size [2]int

	for key, nodeInstance := range nodes {
		outFileNodeVal.WriteString(key + " " + strconv.FormatFloat(nodeInstance.GetVoltage(), 'e', 8, 64) + "\n")
		splitedKey := strings.Split(key, "_")

		newX := utils.ParseInt(splitedKey[len(splitedKey)-2])
		newY := utils.ParseInt(splitedKey[len(splitedKey)-1])

		if newX > size[0] {
			size[0] = newX
		}
		if newY > size[1] {
			size[1] = newY
		}

		bar.Add(1)
	}

	for key, v := range volage {
		outFileNodeVal.WriteString(key + " " + strconv.FormatFloat(v, 'e', 8, 64) + "\n")
		bar.Add(1)
	}

	outFileNodeVal.Close()

	bar = prettier.DefaultBar(len(res), "Calculating IRDrop...")
	irdrop := make(map[string]float64)

	var maxDrop float64
	var minDrop float64 = math.Inf(1)

	for key, r := range res {
		var firstVal float64
		var secondVal float64

		splitedR := strings.Split(r, " ")

		if entryNode, foundFirst := nodes[splitedR[1]]; foundFirst {
			firstVal = entryNode.GetVoltage()
		}

		if entryNode, foundSecond := nodes[splitedR[2]]; foundSecond {
			secondVal = entryNode.GetVoltage()
		}

		drop := math.Abs(firstVal - secondVal)
		irdrop[r+" "+key] = drop

		if drop > maxDrop {
			maxDrop = drop
		} else if drop < minDrop {
			minDrop = drop
		}

		bar.Add(1)
	}

	keys := make([]string, 0, len(irdrop))

	for key := range irdrop {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	outFileResIRDrop := utils.CreateFile(resultsPath + "/" + projName + "/" + projName + ".plot")

	outFileResIRDrop.WriteString("# Proj_Name " + filePath + "\n")
	outFileResIRDrop.WriteString("# Size " + strconv.FormatInt(int64(size[0]), 10) + " " + strconv.FormatInt(int64(size[1]), 10) + "\n")
	outFileResIRDrop.WriteString("# Max_Val " + strconv.FormatFloat(maxDrop, 'e', 8, 64) + "\n")
	outFileResIRDrop.WriteString("# Min_Val " + strconv.FormatFloat(minDrop, 'e', 8, 64) + "\n")

	for i := 0; i < len(keys); i++ {
		outFileResIRDrop.WriteString(keys[i] + " " + strconv.FormatFloat(irdrop[keys[i]], 'e', 8, 64) + "\n")
		bar.Add(1)
	}

	bar.Close()
	outFileResIRDrop.Close()
}
