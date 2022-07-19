package model

import (
	"fmt"
	"pgsolver/pkg/prettier"
)

// Zeidel model for solving system of equations.
// Require max steps for modeling and accurasy of method.
type Model struct {
	voltage  map[string]float64
	current  map[string]float64
	nodes    map[string]*Node
	maxSteps int
	e        float64
}

// Craete instance of Model
func NewModel(voltage map[string]float64, current map[string]float64, nodes map[string]*Node, maxSteps int, e float64) *Model {
	model := &Model{voltage: voltage, current: current, nodes: nodes, maxSteps: maxSteps, e: e}
	return model
}

// Preapare model of nodes.
// Connect all nodes with each others.
// Calculate inital state for each modeling node.
func (model *Model) Prepare() {
	bar := prettier.DefaultBar(len(model.nodes), "Preparing the model...")

	// Replace node names by node instance also voltage and current sources
	for key, nodeInstance := range model.nodes {
		if _, found := model.voltage[key]; !found {
			if entryCurrent, found := model.current[key]; found {
				nodeInstance.i += entryCurrent
			}

			if len(nodeInstance.viases) == 0 {
				for i := 0; i < len(nodeInstance.connectedNodes); i++ {
					if entryNode, ok := nodeInstance.connectedNodes[i].(string); ok {
						if entryV, found := model.voltage[entryNode]; found {
							nodeInstance.connectedNodes[i] = entryV
						} else {
							if entryNode, found := model.nodes[entryNode]; found {
								nodeInstance.connectedNodes[i] = entryNode
							}
						}
					}
				}
			} else {
				for i := 0; i < len(nodeInstance.viases); i++ {
					if entryNode, found := model.nodes[nodeInstance.viases[i]]; found {
						nodeInstance.viasesNodes = append(nodeInstance.viasesNodes, entryNode.connectedNodes...)
						nodeInstance.viasesRes = append(nodeInstance.viasesRes, entryNode.connectedRes...)
					}

					if entryCurrent, found := model.current[nodeInstance.viases[i]]; found {
						nodeInstance.i += entryCurrent
					}
				}

				for i := 0; i < len(nodeInstance.connectedNodes); i++ {
					if entryNode, ok := nodeInstance.connectedNodes[i].(string); ok {
						if entryV, found := model.voltage[entryNode]; found {
							nodeInstance.connectedNodes[i] = entryV
						} else {
							if entryNode, found := model.nodes[entryNode]; found {
								nodeInstance.connectedNodes[i] = entryNode
							}
						}
					}
				}

				for i := 0; i < len(nodeInstance.viasesNodes); i++ {
					if entryNode, ok := nodeInstance.viasesNodes[i].(string); ok {
						if entryV, found := model.voltage[entryNode]; found {
							nodeInstance.viasesNodes[i] = entryV
						} else {
							if entryNode, found := model.nodes[entryNode]; found {
								nodeInstance.viasesNodes[i] = entryNode
							}
						}
					}
				}

			}
		} else {
			delete(model.nodes, key)
		}

		bar.Add(1)
	}

	bar.Close()
	fmt.Println()
}

// Initializing model.
// Set default values, stack nodes and res.
func (model *Model) Init() {
	bar := prettier.DefaultBar(len(model.nodes), "Initializing the model...")

	for _, nodeInstance := range model.nodes {
		nodeInstance.Init()
		bar.Add(1)
	}

	bar.Close()
	fmt.Println()
}

// Modeling system of nodes.
// Modeling last until max steps achived.
func (model *Model) Modeling() {
	bar := prettier.DefaultBar(model.maxSteps, "Solving the model...")
	total := 0

	for i := 0; i < model.maxSteps; i++ {
		solvedNodes := 0

		for _, nodeInstance := range model.nodes {
			solvedNodes += nodeInstance.Step(model.e)
		}

		if solvedNodes != len(model.nodes) {
			bar.Add(1)
			total += 1
		} else {
			break
		}
	}

	bar.Close()

	if total == model.maxSteps {
		fmt.Printf("\n\nAccuracy has not been achieved.")
	} else {
		fmt.Printf("\n\nAccuracy has been achieved with steps count: %d\n", total+1)
	}

	fmt.Println()
}
