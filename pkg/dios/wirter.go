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
func WriteLogs(nodes map[string]*model.Node, res map[string]string, fileName string) {
	splitedPath := strings.Split(fileName, "/")
	projName := strings.Split(splitedPath[len(splitedPath)-1], ".")[0]

	utils.CreateFolder("./results")
	utils.CreateFolder("./results/" + projName)

	bar := prettier.DefaultBar(len(nodes)+2*len(res), "Writing logs...")
	outFileNodeVal := utils.CreateFile("./results/" + projName + "/" + projName + ".solution")

	var size [2]int

	for key, nodeInstance := range nodes {
		outFileNodeVal.WriteString(key + " " + strconv.FormatFloat(nodeInstance.V, 'e', 8, 64) + "\n")
		splitedKey := strings.Split(key, "_")

		newX := utils.ParseInt(splitedKey[1])
		newY := utils.ParseInt(splitedKey[2])

		if newX > size[0] {
			size[0] = newX
		}
		if newY > size[1] {
			size[1] = newY
		}

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
			firstVal = entryNode.V
		}

		if entryNode, foundSecond := nodes[splitedR[2]]; foundSecond {
			secondVal = entryNode.V
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

	outFileResIRDrop := utils.CreateFile("./results/" + projName + "/" + projName + ".plot")

	outFileResIRDrop.WriteString("# Proj_Name " + fileName + "\n")
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
