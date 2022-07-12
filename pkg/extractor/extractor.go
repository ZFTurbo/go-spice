package extractor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"pgsolver/pkg/node"
	"pgsolver/pkg/prettier"
	"strconv"
	"strings"
)

// Data extractor for ibm spice formate files.
type Extractor struct {
	fileName string
}

// Craetes new instance of Extractor.
func NewExtractor(fileName string) *Extractor {
	extractor := &Extractor{fileName: fileName}
	return extractor
}

// Extract all data from spice net list.
// Return voltage, current, node maps.
func (extractor *Extractor) Extract() (map[string]float64, map[string]float64, map[string]*node.Node, error) {

	consPrettier := prettier.NewPrettier()
	voltage := make(map[string]float64)
	current := make(map[string]float64)
	nodes := make(map[string]*node.Node)

	file, err := os.Open(extractor.fileName)

	if err == nil {
		if fileStat, err := file.Stat(); err == nil {
			bar := consPrettier.DefaultBar(int(fileStat.Size()), "Extraction from input file...")
			scanner := bufio.NewScanner(file) // The file scanner

			// Reading input file
			for scanner.Scan() {
				line := scanner.Text()
				splitedLine := strings.Split(line, " ")
				// Find volage source
				if line[0] == 'v' {
					if entryV, err := strconv.ParseFloat(splitedLine[len(splitedLine)-1], 64); err == nil {
						voltage[splitedLine[1]] = entryV
					} else {
						bar.Close()
						fmt.Println()
						return nil, nil, nil, errors.New("Add voltage error.\n" + err.Error())
					}
				}

				// Find current source
				if line[0] == 'i' {
					if entryI, err := strconv.ParseFloat(splitedLine[len(splitedLine)-2], 64); err == nil {
						// By default use as node name as first connection of current source, not ground connection
						if strings.Contains(splitedLine[0], "_v") {
							if entryCurrent, found := current[splitedLine[1]]; found {
								entryCurrent += entryI
								current[splitedLine[1]] = entryCurrent
							} else {
								current[splitedLine[1]] = entryI
							}
						} else {
							// Due to IBM format if current source direction to ground change current dirction, not ground connection
							if strings.Contains(splitedLine[0], "_g") {
								if entryCurrent, found := current[splitedLine[2]]; found {
									entryCurrent += -entryI
									current[splitedLine[2]] = entryCurrent
								} else {
									current[splitedLine[2]] = -entryI
								}
							}
						}
					} else {
						bar.Close()
						fmt.Println()
						return nil, nil, nil, errors.New("Add current ground error.\n" + err.Error())
					}
				}

				if (line[0] == 'r' || line[0] == 'R' || line[0] == 'V') && splitedLine[len(splitedLine)-2] != "0" {
					if entryRes, err := strconv.ParseFloat(splitedLine[len(splitedLine)-1], 64); err == nil {
						// Check if resistor is via.
						// If resistance value is too small.
						if entryRes != 0.0 {
							if entryNode, found := nodes[splitedLine[1]]; found {
								entryNode.ConnectedNodes = append(entryNode.ConnectedNodes, splitedLine[2])
								entryNode.ConnectedRes = append(entryNode.ConnectedRes, entryRes)
							} else {
								nodes[splitedLine[1]] = node.NewNode(splitedLine[1], splitedLine[2], entryRes)
							}

							if entryNode, found := nodes[splitedLine[2]]; found {
								entryNode.ConnectedNodes = append(entryNode.ConnectedNodes, splitedLine[1])
								entryNode.ConnectedRes = append(entryNode.ConnectedRes, entryRes)
							} else {
								nodes[splitedLine[2]] = node.NewNode(splitedLine[2], splitedLine[1], entryRes)
							}
						} else {
							if entryVia, found := nodes[splitedLine[1]]; found {
								entryVia.Viases = append(entryVia.Viases, splitedLine[2])
							} else {
								nodes[splitedLine[1]] = node.NewVia(splitedLine[1], splitedLine[2])
							}

							if entryVia, found := nodes[splitedLine[2]]; found {
								entryVia.Viases = append(entryVia.Viases, splitedLine[1])
							} else {
								nodes[splitedLine[2]] = node.NewVia(splitedLine[2], splitedLine[1])
							}
						}
					} else {
						bar.Close()
						fmt.Println()
						return nil, nil, nil, errors.New("Create node. Resistance value error.\n" + err.Error())
					}

					bar.Add(1)
				}
			}

			bar.Close()
			fmt.Println()
			return voltage, current, nodes, nil
		} else {
			return nil, nil, nil, errors.New("Can't get size of file!\n" + err.Error())
		}
	} else {
		return nil, nil, nil, errors.New("Can't open file!\n" + err.Error())
	}
}
