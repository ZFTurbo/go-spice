package model

import (
	"fmt"

	"pgsolver/pkg/prettier"
)

// Zeidel model for solving system of equations.
// Require max steps for modeling and accurasy of method.
type Model struct {
	nodes    map[string]*Node
	maxSteps int
	e        float64
}

// Craete instance of Model
func NewModel(nodes map[string]*Node, maxSteps int, e float64) *Model {
	model := &Model{}
	model.nodes = nodes
	model.maxSteps = maxSteps
	model.e = e

	return model
}

// Initializing model.
// Set default values, stack nodes and res.
func (model *Model) initElements() {
	bar := prettier.DefaultBar(len(model.nodes), "Initializing the model elements...")

	for _, nodeInstance := range model.nodes {
		nodeInstance.Init()
		bar.Add(1)
	}

	bar.Close()
	fmt.Println()
}

// DC anlysis modeling
func DCModeling(model *Model) {
	bar := prettier.DefaultBar(model.maxSteps, "Solving the model in dc analysis...")
	total := 0

	model.initElements()

	for i := 0; i < model.maxSteps; i++ {
		solvedNodes := 0

		for _, nodeInstance := range model.nodes {
			solvedNodes += nodeInstance.StepDC(model.e)
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
		fmt.Printf("\n\nAccuracy has not been achieved.\n")
	} else {
		fmt.Printf("\n\nAccuracy has been achieved with steps count: %d\n", total+1)
	}

	fmt.Println()
}
