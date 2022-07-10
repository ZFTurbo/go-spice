package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"pgsolver/pkg/node"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {
	timerStart := time.Now() // Exeqution timer

	inFilePath := flag.String("f", "../data/ibmpg1.spice", "Iput file")      // Path to input file
	outFilePath := flag.String("o", "../out/ibmpg1.solution", "Output file") // Path to output file
	e := flag.Float64("p", 1e-8, "Precision of modeling")                    // Accuracy of the Zeidele method
	maxSteps := flag.Int("ms", 100000, "Max count of steps in modeling")     // Max ammount of step during modeling

	voltage := make(map[string]float64)  // Map of all voltage sources
	current := make(map[string]float64)  // Map of all current sources
	nodes := make(map[string]*node.Node) // Map of all nodes

	bar := progressbar.Default(int64(*maxSteps)) // Progress bar

	flag.Parse()
	fmt.Println("\nExtracting data from file...")

	// Open file
	inFile, err := os.Open(*inFilePath)

	if err != nil {
		log.Fatal("Error in input file.\n", err)
	}

	scanner := bufio.NewScanner(inFile) // The file scanner

	// Reading input file
	for scanner.Scan() {
		line := scanner.Text()
		splitedLine := strings.Split(line, " ")
		// Find volage source
		if line[0] == 'v' {
			if entryV, err := strconv.ParseFloat(splitedLine[len(splitedLine)-1], 64); err == nil {
				voltage[splitedLine[1]] = entryV
			} else {
				log.Fatal("Add voltage error: ", err)
			}
		}

		// Find current source
		if line[0] == 'i' {
			if entryI, err := strconv.ParseFloat(splitedLine[len(splitedLine)-2], 64); err == nil {
				if strings.Contains(splitedLine[0], "_g") {
					if entryCurrent, found := current[splitedLine[2]]; found {
						entryCurrent += -entryI
						current[splitedLine[2]] = entryCurrent
					} else {
						current[splitedLine[2]] = -entryI
					}
				}

				if strings.Contains(splitedLine[0], "_v") {
					if entryCurrent, found := current[splitedLine[1]]; found {
						entryCurrent += entryI
						current[splitedLine[1]] = entryCurrent
					} else {
						current[splitedLine[1]] = entryI
					}
				}
			} else {
				log.Fatal("Add current ground error. ", err)
			}
		}

		if (line[0] == 'r' || line[0] == 'R' || line[0] == 'V') && splitedLine[len(splitedLine)-2] != "0" {
			if entryRes, err := strconv.ParseFloat(splitedLine[len(splitedLine)-1], 64); err == nil {
				// Check if resistor is via.
				// If resistance value is too small.
				if entryRes > 1e-06 {
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
				log.Fatal("Create node. Resistance value error. ", err)
			}
		}
	}

	inFile.Close()

	fmt.Println("\nPreparing node based model...")

	// Replace node names by node instance also voltage and current sources
	for _, nodeInstance := range nodes {
		if len(nodeInstance.Viases) == 0 {
			if entryCurrent, found := current[nodeInstance.Name]; found {
				nodeInstance.I += entryCurrent
			}
			for i := 0; i < len(nodeInstance.ConnectedNodes); i++ {
				if entryNode, ok := nodeInstance.ConnectedNodes[i].(string); ok {
					if entryV, found := voltage[entryNode]; found {
						nodeInstance.ConnectedNodes[i] = entryV
					} else {
						if entryNode, found := nodes[entryNode]; found {
							nodeInstance.ConnectedNodes[i] = entryNode
						}
					}
				}
			}
		} else {
			if entryCurrent, found := current[nodeInstance.Name]; found {
				nodeInstance.I += entryCurrent
			}

			for i := 0; i < len(nodeInstance.Viases); i++ {
				if entryNode, found := nodes[nodeInstance.Viases[i]]; found {
					nodeInstance.ViasesNodes = append(nodeInstance.ViasesNodes, entryNode.ConnectedNodes...)
					nodeInstance.ViasesRes = append(nodeInstance.ViasesRes, entryNode.ConnectedRes...)
				}

				if entryCurrent, found := current[nodeInstance.Viases[i]]; found {
					nodeInstance.I += entryCurrent
				}
			}

			for i := 0; i < len(nodeInstance.ConnectedNodes); i++ {
				if entryNode, ok := nodeInstance.ConnectedNodes[i].(string); ok {
					if entryV, found := voltage[entryNode]; found {
						nodeInstance.ConnectedNodes[i] = entryV
					} else {
						if entryNode, found := nodes[entryNode]; found {
							nodeInstance.ConnectedNodes[i] = entryNode
						}
					}
				}
			}

			for i := 0; i < len(nodeInstance.ViasesNodes); i++ {
				if entryNode, ok := nodeInstance.ViasesNodes[i].(string); ok {
					if entryV, found := voltage[entryNode]; found {
						nodeInstance.ViasesNodes[i] = entryV
					} else {
						if entryNode, found := nodes[entryNode]; found {
							nodeInstance.ViasesNodes[i] = entryNode
						}
					}
				}
			}

		}

	}

	fmt.Println("\nSoving node model...")
	fmt.Println()

	// Init all nodes
	for _, nodeInstance := range nodes {
		nodeInstance.Init()
	}

	// Modeling system of nodes
	for i := 0; i < *maxSteps; i++ {
		solvedNodes := 0

		for _, nodeInstance := range nodes {
			solvedNodes += nodeInstance.Step(*e)
		}

		if solvedNodes == len(nodes) {
			fmt.Println()
			fmt.Println("\nSteps cout: ", i+1)
			fmt.Println()
			break
		}

		bar.Add(1)
	}

	fmt.Println()
	fmt.Println("\nSteps cout: ", *maxSteps)
	fmt.Println()
	fmt.Println("\nWriting results...")

	// Log out data
	outFile, err := os.Create(*outFilePath)
	if err != nil {
		log.Fatal("Create output file error.\n",err)
	} else {
		for _, nodeInstance := range nodes {
			outFile.WriteString(nodeInstance.Name + " ")
			outFile.WriteString(strconv.FormatFloat(nodeInstance.V, 'e', 8, 64))
			outFile.WriteString("\n")
		}
	}

	outFile.Close()

	fmt.Println("\nEnd of modeling!")
	fmt.Println()
	fmt.Println("Exeqution time in seconds: ", time.Now().Sub(timerStart).Seconds())
}
