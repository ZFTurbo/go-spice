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
		if entryCurrent, found := model.current[key]; found {
			nodeInstance.I += entryCurrent
		}

		if len(nodeInstance.Viases) == 0 {
			for i := 0; i < len(nodeInstance.ConnectedNodes); i++ {
				if entryNode, ok := nodeInstance.ConnectedNodes[i].(string); ok {
					if entryV, found := model.voltage[entryNode]; found {
						nodeInstance.ConnectedNodes[i] = entryV
					} else {
						if entryNode, found := model.nodes[entryNode]; found {
							nodeInstance.ConnectedNodes[i] = entryNode
						}
					}
				}
			}
		} else {
			for i := 0; i < len(nodeInstance.Viases); i++ {
				if entryNode, found := model.nodes[nodeInstance.Viases[i]]; found {
					nodeInstance.ViasesNodes = append(nodeInstance.ViasesNodes, entryNode.ConnectedNodes...)
					nodeInstance.ViasesRes = append(nodeInstance.ViasesRes, entryNode.ConnectedRes...)
				}

				if entryCurrent, found := model.current[nodeInstance.Viases[i]]; found {
					nodeInstance.I += entryCurrent
				}
			}

			for i := 0; i < len(nodeInstance.ConnectedNodes); i++ {
				if entryNode, ok := nodeInstance.ConnectedNodes[i].(string); ok {
					if entryV, found := model.voltage[entryNode]; found {
						nodeInstance.ConnectedNodes[i] = entryV
					} else {
						if entryNode, found := model.nodes[entryNode]; found {
							nodeInstance.ConnectedNodes[i] = entryNode
						}
					}
				}
			}

			for i := 0; i < len(nodeInstance.ViasesNodes); i++ {
				if entryNode, ok := nodeInstance.ViasesNodes[i].(string); ok {
					if entryV, found := model.voltage[entryNode]; found {
						nodeInstance.ViasesNodes[i] = entryV
					} else {
						if entryNode, found := model.nodes[entryNode]; found {
							nodeInstance.ViasesNodes[i] = entryNode
						}
					}
				}
			}

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
	p := int(0.5 * float64(len(model.nodes)))
	total := 0

	for i := 0; i < model.maxSteps; i++ {
		solvedNodes := 0

		for _, nodeInstance := range model.nodes {
			solvedNodes += nodeInstance.Step(model.e)
		}

		if solvedNodes < p {
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
