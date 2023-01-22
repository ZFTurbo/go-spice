package dios

import (
	"bufio"
	"fmt"
	"strings"

	"pgsolver/pkg/model"
	"pgsolver/pkg/prettier"
	"pgsolver/pkg/utils"
)

/*
Extract all data from spice net list.
Return voltage, current, node maps.
*/
func Extract(fileName string) map[string]*model.Node {
	nodes := make(map[string]*model.Node)

	file := utils.OpenFile(fileName)

	if fileStat, err := file.Stat(); err == nil {
		bar := prettier.DefaultBar(int(fileStat.Size()), "Extraction from input file...")
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := strings.Trim(strings.ToLower(scanner.Text()), " ")

			if line[0] != '.' && line[0] != '*' {
				splitedLine := strings.Fields(line)
				lastElement := len(splitedLine) - 1
				lineValue := utils.ParseFloat(splitedLine[lastElement])

				if _, found := nodes[splitedLine[1]]; !found {
					nodes[splitedLine[1]] = model.NewNode(splitedLine[1])
				}

				if _, found := nodes[splitedLine[2]]; !found {
					nodes[splitedLine[2]] = model.NewNode(splitedLine[2])
				}

				if splitedLine[1] != "0" {
					nodes[splitedLine[2]].AddNode(nodes[splitedLine[1]])
				}

				if splitedLine[2] != "0" {
					nodes[splitedLine[1]].AddNode(nodes[splitedLine[2]])
				}

				switch line[0] {
				case 'v':
					nodes[splitedLine[1]].SetModeling(false)
					nodes[splitedLine[2]].SetModeling(false)

					if splitedLine[1] != "0" {
						nodes[splitedLine[1]].SetVoltage(lineValue)
					}

					if splitedLine[2] != "0" {
						nodes[splitedLine[2]].SetVoltage(-lineValue)
					}

				case 'i':
					if splitedLine[1] != "0" {
						nodes[splitedLine[1]].AddCurrentSource(model.NewCurrent(splitedLine[1], -lineValue))
					} else {
						nodes[splitedLine[1]].SetModeling(false)
					}

					if splitedLine[2] != "0" {
						nodes[splitedLine[2]].AddCurrentSource(model.NewCurrent(splitedLine[2], lineValue))
					} else {
						nodes[splitedLine[2]].SetModeling(false)
					}

				case 'r':
					if lineValue != 0.0 {
						if splitedLine[1] != "0" {
							nodes[splitedLine[1]].AddResistor(lineValue)
						} else {
							nodes[splitedLine[1]].SetModeling(false)
						}

						if splitedLine[2] != "0" {
							nodes[splitedLine[2]].AddResistor(lineValue)
						} else {
							nodes[splitedLine[2]].SetModeling(false)
						}

					}
				}
			}

			bar.Add(1)
		}

		bar.Close()
		fmt.Println()
		fmt.Println()

		return nodes
	} else {
		prettier.Error("Can't get file stats.", err)
		return nil
	}
}
