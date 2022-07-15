package dios

import (
	"math"
	"os"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"sort"
	"strconv"
	"strings"
)

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Writes solution in file
func WriteLogs(nodes map[string]*model.Node, res map[string]string, fileName string) {

	splitedPath := strings.Split(fileName, "/")
	projName := strings.Split(splitedPath[len(splitedPath)-1], ".")[0]

	if ok, err := exists("./results"); err == nil {
		if !ok {
			err = os.Mkdir("./results", os.ModePerm)
			if err != nil {
				prettier.Error("Error has been occure while creating results folder.", err)
			}
		}
	} else {
		prettier.Error("Error has been occure while folder exists.", err)
	}

	if ok, err := exists("./results/" + projName); err == nil {
		if !ok {
			err = os.Mkdir("./results/"+projName, os.ModePerm)
			if err != nil {
				prettier.Error("Error has been occure while creating proj folder in results folder.", err)
			}
		}
	} else {
		prettier.Error("Error has been occure while folder exists.", err)
	}

	bar := prettier.DefaultBar(len(nodes)+2*len(res), "Writing logs...")
	var size [2]int

	outFileNodeVal, err := os.Create("./results/" + projName + "/" + projName + ".solution")

	if err != nil {
		prettier.Error("Create output file error.", err)
	} else {
		for key, nodeInstance := range nodes {
			outFileNodeVal.WriteString(key + " " + strconv.FormatFloat(nodeInstance.V, 'e', 8, 64) + "\n")
			splitedKey := strings.Split(key, "_")
			if entryX, err := strconv.ParseInt(splitedKey[1], 10, 0); err == nil {
				if int(entryX) > size[0] {
					size[0] = int(entryX)
				}
			}
			if entryY, err := strconv.ParseInt(splitedKey[2], 10, 0); err == nil {
				if int(entryY) > size[1] {
					size[1] = int(entryY)
				}
			}
			bar.Add(1)
		}
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

	outFileResIRDrop, err := os.Create("./results/" + projName + "/" + projName + ".plot")

	if err != nil {
		prettier.Error("Create output file error.", err)
	} else {
		outFileResIRDrop.WriteString("# Proj_Name " + fileName + "\n")
		outFileResIRDrop.WriteString("# Size " + strconv.FormatInt(int64(size[0]), 10) + " " + strconv.FormatInt(int64(size[1]), 10) + "\n")
		outFileResIRDrop.WriteString("# Max_Val " + strconv.FormatFloat(maxDrop, 'e', 8, 64) + "\n")
		outFileResIRDrop.WriteString("# Min_Val " + strconv.FormatFloat(minDrop, 'e', 8, 64) + "\n")

		for i := 0; i < len(keys); i++ {
			outFileResIRDrop.WriteString(keys[i] + " " + strconv.FormatFloat(irdrop[keys[i]], 'e', 8, 64) + "\n")
			bar.Add(1)
		}
	}

	bar.Close()
	outFileResIRDrop.Close()
}
