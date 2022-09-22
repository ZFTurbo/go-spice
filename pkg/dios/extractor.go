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
func Extract(fileName string) (*model.Analysis, map[string]string, map[string]float64, map[string][]*model.Current, map[string][]*model.Capaster, map[string][]*model.Inductance, map[string]*model.Node) {
	re := regexp.MustCompile(`-?[\d.]+(?:e-?\d+)?`)
	analysis := model.NewAnalysis("DC", 0, 0)
	res := make(map[string]string)
	voltage := make(map[string]float64)
	current := make(map[string][]*model.Current)
	inductance := make(map[string][]*model.Inductance)
	capasters := make(map[string][]*model.Capaster)
	nodes := make(map[string]*model.Node)

	file := utils.OpenFile(fileName)

	if fileStat, err := file.Stat(); err == nil {
		bar := prettier.DefaultBar(int(fileStat.Size()), "Extraction from input file...")
		scanner := bufio.NewScanner(file) // The file scanner
		resType := "VDD"

		//Add ground
		voltage["0"] = 0

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
				// TODO: Size down this block of code
				if entryNode, found := nodes[splitedLine[1]]; found {
					indc := model.NewInductance(splitedLine[2], utils.ParseFloat(splitedLine[lastElement]))
					entryNode.AddInductance(indc)
				} else {
					indc := model.NewInductance(splitedLine[2], utils.ParseFloat(splitedLine[lastElement]))
					nodes[splitedLine[1]] = model.NewNodeRaw(splitedLine[1])
					nodes[splitedLine[1]].AddInductance(indc)
				}

				if entryNode, found := nodes[splitedLine[2]]; found {
					indc := model.NewInductance(splitedLine[1], utils.ParseFloat(splitedLine[lastElement]))
					entryNode.AddInductance(indc)
				} else {
					indc := model.NewInductance(splitedLine[1], utils.ParseFloat(splitedLine[lastElement]))
					nodes[splitedLine[2]] = model.NewNodeRaw(splitedLine[2])
					nodes[splitedLine[2]].AddInductance(indc)
				}
			case 'c':
				// TODO: Size down this block of code
				if entryNode, found := nodes[splitedLine[1]]; found {
					cap := model.NewCapaster(splitedLine[2], utils.ParseFloat(splitedLine[lastElement]))
					entryNode.AddCapaster(cap)
				} else {
					cap := model.NewCapaster(splitedLine[2], utils.ParseFloat(splitedLine[lastElement]))
					nodes[splitedLine[1]] = model.NewNodeRaw(splitedLine[1])
					nodes[splitedLine[1]].AddCapaster(cap)
				}

				if entryNode, found := nodes[splitedLine[2]]; found {
					cap := model.NewCapaster(splitedLine[1], -utils.ParseFloat(splitedLine[lastElement]))
					entryNode.AddCapaster(cap)
				} else {
					cap := model.NewCapaster(splitedLine[1], -utils.ParseFloat(splitedLine[lastElement]))
					nodes[splitedLine[2]] = model.NewNodeRaw(splitedLine[2])
					nodes[splitedLine[2]].AddCapaster(cap)
				}
			case 'v':
				voltage[splitedLine[1]] = utils.ParseFloat(splitedLine[lastElement])
			case 'i':
				iValue := utils.ParseFloat(splitedLine[3])

				// By default use as node name as first connection of current source, not ground connection
				if strings.Contains(splitedLine[0], "_v") {
					if entryCurrent, found := current[splitedLine[1]]; found {
						entryCurrent = append(entryCurrent, model.NewCurrent(splitedLine[1], iValue))
						current[splitedLine[1]] = entryCurrent
					} else {
						if len(splitedLine) < 5 {
							current[splitedLine[1]] = append(current[splitedLine[1]], model.NewCurrent(splitedLine[1], iValue))
						} else {
							// TODO: Replace this by function.
							if strings.Contains(line, "pulse") {
								pulseSplited := re.FindAllString(strings.Split(line, "pulse")[1], -1)
								max := utils.ParseFloat(pulseSplited[1])
								td := utils.ParseFloat(pulseSplited[2])
								tr := utils.ParseFloat(pulseSplited[3])
								tf := utils.ParseFloat(pulseSplited[4])
								pw := utils.ParseFloat(pulseSplited[5])
								per := utils.ParseFloat(pulseSplited[6])

								current[splitedLine[1]] = append(current[splitedLine[1]], model.NewCurrentPulse(splitedLine[1], iValue, max, td, tr, tf, pw, per))
							}
						}
					}
				} else {
					// Due to IBM format if current source direction to ground change current dirction, not ground connection
					if strings.Contains(splitedLine[0], "_g") {
						if entryCurrent, found := current[splitedLine[2]]; found {
							entryCurrent = append(entryCurrent, model.NewCurrent(splitedLine[2], -iValue))
							current[splitedLine[2]] = entryCurrent
						} else {
							if len(splitedLine) < 5 {
								current[splitedLine[2]] = append(current[splitedLine[2]], model.NewCurrent(splitedLine[2], -iValue))
							} else {
								// TODO: Replace this by function.
								if strings.Contains(line, "pulse") {
									pulseSplited := re.FindAllString(strings.Split(line, "pulse")[1], -1)
									max := utils.ParseFloat(pulseSplited[1])
									td := utils.ParseFloat(pulseSplited[2])
									tr := utils.ParseFloat(pulseSplited[3])
									tf := utils.ParseFloat(pulseSplited[4])
									pw := utils.ParseFloat(pulseSplited[5])
									per := utils.ParseFloat(pulseSplited[6])

									current[splitedLine[2]] = append(current[splitedLine[2]], model.NewCurrentPulse(splitedLine[2], -iValue, -max, td, tr, tf, pw, per))
								}
							}
						}
					}
				}
			case 'r', 'R', 'V':
				resVal := utils.ParseFloat(splitedLine[lastElement])

				// FIXME: Replace with regexp
				if (line[0] == 'R' || line[0] == 'r') && splitedLine[lastElement-1] != "0" && splitedLine[1] != "0" && splitedLine[1][1] == splitedLine[2][1] {
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
			case '.':
				if strings.Contains(line, ".tran") {
					analysis.SetName("TR")
					analysis.SetFullTime(utils.ParseFloat(splitedLine[2]))
					analysis.SetTimeStep(utils.ParseFloat(splitedLine[1]))
				}
			}

			bar.Add(1)
		}

		bar.Close()
		fmt.Println()
		fmt.Println()

		return analysis, res, voltage, current, capasters, inductance, nodes
	} else {
		prettier.Error("Can't get file stats.", err)
		return nil, nil, nil, nil, nil, nil, nil
	}
}
