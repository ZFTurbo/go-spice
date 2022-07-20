package dios

import (
	"bufio"
	"fmt"
	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"pgsolver/pkg/utils"
	"regexp"
	"strings"
)

// Extract all data from spice net list.
// Return voltage, current, node maps.
func Extract(fileName string) (map[string]string, map[string]float64, map[string]*model.Current, map[string]*model.Node) {
	re := regexp.MustCompile(`-?[\d.]+(?:e-?\d+)?`)
	res := make(map[string]string)
	voltage := make(map[string]float64)
	current := make(map[string]*model.Current)
	inductance := make(map[string]float64)
	capasters := make(map[string]float64)
	nodes := make(map[string]*model.Node)

	file := utils.OpenFile(fileName)

	if fileStat, err := file.Stat(); err == nil {
		bar := prettier.DefaultBar(int(fileStat.Size()), "Extraction from input file...")
		scanner := bufio.NewScanner(file) // The file scanner
		resType := "VDD"

		// Reading input file
		for scanner.Scan() {
			line := scanner.Text()
			splitedLine := strings.Fields(line)
			lastElement := len(splitedLine) - 1

			switch line[0] {
			case '*':
				if strings.Contains(line, "VDD") {
					resType = "VDD"
				}
				if strings.Contains(line, "GND") {
					resType = "GND"
				}
			case 'l':
				inductance[splitedLine[1]] = utils.ParseFloat(splitedLine[lastElement])
			case 'c':
				capasters[splitedLine[1]] = utils.ParseFloat(splitedLine[lastElement])
			case 'v':
				voltage[splitedLine[1]] = utils.ParseFloat(splitedLine[lastElement])
			case 'i':
				iValue := utils.ParseFloat(splitedLine[3])

				// By default use as node name as first connection of current source, not ground connection
				if strings.Contains(splitedLine[0], "_v") {
					if entryCurrent, found := current[splitedLine[1]]; found {
						entryCurrent.AddVal(iValue)
					} else {
						if len(splitedLine) < 5 {
							current[splitedLine[1]] = model.NewCurrent(splitedLine[1], iValue)
						} else {
							if strings.Contains(line, "pulse") {
								pulseSplited := re.FindAllString(strings.Split(line, "pulse")[1], -1)
								min := utils.ParseFloat(pulseSplited[1])
								tr := utils.ParseFloat(pulseSplited[2])
								tf := utils.ParseFloat(pulseSplited[3])
								pw := utils.ParseFloat(pulseSplited[4])
								dl := utils.ParseFloat(pulseSplited[5])

								current[splitedLine[1]] = model.NewCurrentPulse(splitedLine[1], iValue, min, tr, tf, pw, dl)
							}
						}
					}
				} else {
					// Due to IBM format if current source direction to ground change current dirction, not ground connection
					if strings.Contains(splitedLine[0], "_g") {
						if entryCurrent, found := current[splitedLine[2]]; found {
							entryCurrent.AddVal(-iValue)
						} else {
							if len(splitedLine) < 5 {
								current[splitedLine[2]] = model.NewCurrent(splitedLine[2], -iValue)
							} else {
								if strings.Contains(line, "pulse") {
									pulseSplited := re.FindAllString(strings.Split(line, "pulse")[1], -1)
									min := utils.ParseFloat(pulseSplited[1])
									tr := utils.ParseFloat(pulseSplited[2])
									tf := utils.ParseFloat(pulseSplited[3])
									pw := utils.ParseFloat(pulseSplited[4])
									dl := utils.ParseFloat(pulseSplited[5])

									current[splitedLine[2]] = model.NewCurrentPulse(splitedLine[2], -iValue, min, tr, tf, pw, dl)
								}
							}
						}
					}
				}
			case 'r', 'R', 'V':
				resVal := utils.ParseFloat(splitedLine[lastElement])

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
