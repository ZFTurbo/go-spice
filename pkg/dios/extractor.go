package dios

import (
	"bufio"
	"fmt"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"pgsolver/pkg/utils"
	"strings"
)

// Extract all data from spice net list.
// Return voltage, current, node maps.
func Extract(fileName string) (map[string]string, map[string]float64, map[string]float64, map[string]*model.Node) {
	res := make(map[string]string)
	voltage := make(map[string]float64)
	current := make(map[string]float64)
	nodes := make(map[string]*model.Node)

	file := utils.OpenFile(fileName)

	if fileStat, err := file.Stat(); err == nil {
		bar := prettier.DefaultBar(int(fileStat.Size()), "Extraction from input file...")
		scanner := bufio.NewScanner(file) // The file scanner
		resType := "VDD"

		// Reading input file
		for scanner.Scan() {
			line := scanner.Text()
			splitedLine := strings.Split(line, " ")

			// Change type of resisotr gnd or vdd(vpwr)
			if line[0] == '*' && len(splitedLine) > 2 {
				if strings.Contains(splitedLine[2], "VDD") {
					resType = "VDD"
				}
				if strings.Contains(splitedLine[2], "GND") {
					resType = "GND"
				}
			}

			// Find volage source
			if line[0] == 'v' {
				voltage[splitedLine[1]] = utils.ParseFloat(splitedLine[len(splitedLine)-1])
			}

			// Find current source
			if line[0] == 'i' {
				iValue := utils.ParseFloat(splitedLine[len(splitedLine)-2])
				// By default use as node name as first connection of current source, not ground connection
				if strings.Contains(splitedLine[0], "_v") {
					if entryCurrent, found := current[splitedLine[1]]; found {
						entryCurrent += iValue
						current[splitedLine[1]] = entryCurrent
					} else {
						current[splitedLine[1]] = iValue
					}
				} else {
					// Due to IBM format if current source direction to ground change current dirction, not ground connection
					if strings.Contains(splitedLine[0], "_g") {
						if entryCurrent, found := current[splitedLine[2]]; found {
							entryCurrent += -iValue
							current[splitedLine[2]] = entryCurrent
						} else {
							current[splitedLine[2]] = -iValue
						}
					}
				}
			}

			if (line[0] == 'r' || line[0] == 'R' || line[0] == 'V') && splitedLine[len(splitedLine)-2] != "0" {
				resVal := utils.ParseFloat(splitedLine[len(splitedLine)-1])

				if line[0] == 'R' && splitedLine[1][1] == splitedLine[2][1] {
					res[splitedLine[0]] = resType + " " + splitedLine[1] + " " + splitedLine[2]
				}
				// Check if resistor is via.
				// If resistance value is too small.
				if resVal != 0.0 {
					if entryNode, found := nodes[splitedLine[1]]; found {
						entryNode.AddNode(splitedLine[2])
						entryNode.AddRes(resVal)
					} else {
						nodes[splitedLine[1]] = model.NewNode(splitedLine[1], splitedLine[2], resVal)
					}

					if entryNode, found := nodes[splitedLine[2]]; found {
						entryNode.AddNode(splitedLine[1])
						entryNode.AddRes(resVal)
					} else {
						nodes[splitedLine[2]] = model.NewNode(splitedLine[2], splitedLine[1], resVal)
					}
				} else {
					if entryVia, found := nodes[splitedLine[1]]; found {
						entryVia.AddVia(splitedLine[2])
					} else {
						nodes[splitedLine[1]] = model.NewVia(splitedLine[1], splitedLine[2])
					}

					if entryVia, found := nodes[splitedLine[2]]; found {
						entryVia.AddVia(splitedLine[1])
					} else {
						nodes[splitedLine[2]] = model.NewVia(splitedLine[2], splitedLine[1])
					}
				}
			}
			bar.Add(1)
		}
		bar.Close()
		fmt.Println()
		fmt.Println()

		return res, voltage, current, nodes
	} else {
		prettier.Error("Can't get file stats.", err)
		return nil, nil, nil, nil
	}
}
