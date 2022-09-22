package model

import (
	"fmt"
	"pgsolver/pkg/prettier"
	"sort"
)

// Zeidel model for solving system of equations.
// Require max steps for modeling and accurasy of method.
type Model struct {
	voltage    map[string]float64
	current    map[string][]*Current
	inductance map[string][]*Inductance
	capasters  map[string][]*Capaster
	nodes      map[string]*Node
	analysis   *Analysis
	maxSteps   int
	e          float64
}

// Craete instance of Model
func NewModel(voltage map[string]float64, current map[string][]*Current, inductance map[string][]*Inductance, capasters map[string][]*Capaster, nodes map[string]*Node, analysis *Analysis, maxSteps int, e float64) *Model {

	model := &Model{}
	model.voltage = voltage
	model.current = current
	model.inductance = inductance
	model.capasters = capasters
	model.nodes = nodes
	model.analysis = analysis
	model.maxSteps = maxSteps
	model.e = e

	return model
}

// Add viases in shorts with 3 or more nodes
func recursiveAddVia(n *Node, e *Node, model *Model) {
	for _, via := range e.viases {
		if entryNode, found := model.nodes[via]; found && via != n.name {
			if n.AddVia(via) {
				recursiveAddVia(n, entryNode, model)
			}
		}
	}
}

// Disable time depented elenements for dc analysis or
// receiving initial states for transient analysis
func (model *Model) disableTimeDepentedElements() {
	bar := prettier.DefaultBar(len(model.nodes), "Disabling time dependet elements...")

	for _, nodeInstance := range model.nodes {
		for _, indc := range nodeInstance.inductance {
			if val, ok := indc.end.(string); ok {
				nodeInstance.AddVia(val)
			}
			if val, ok := indc.end.(*Node); ok {
				nodeInstance.AddVia(val.name)
			}
		}

		bar.Add(1)
	}

	bar.Close()
	fmt.Println()
}

// Collecting all viases that connected to a node
func (model *Model) collectNodeViases() {
	bar := prettier.DefaultBar(len(model.nodes), "Collecting nodes viases...")

	for key, nodeInstance := range model.nodes {
		for _, via := range nodeInstance.viases {
			if vol, found := model.voltage[via]; found {
				nodeInstance.SetVoltage(vol)
				nodeInstance.SetModeling(false)
			}
		}

		if _, isVoltage := model.voltage[key]; !isVoltage {
			nodeViases := nodeInstance.viases

			for _, via := range nodeViases {
				if entryNode, found := model.nodes[via]; found {
					recursiveAddVia(nodeInstance, entryNode, model)
				}
			}
		}

		bar.Add(1)
	}

	bar.Close()
	fmt.Println()
}

// Preapare model of nodes.
// Connect all nodes with each others.
// Calculate inital state for each modeling node.
func (model *Model) prepareElements() {
	bar := prettier.DefaultBar(len(model.nodes), "Preparing the model elements...")

	// Replace node names by node instance also voltage and current sources
	for key, nodeInstance := range model.nodes {
		if _, isVoltage := model.voltage[key]; !isVoltage {
			if entryCurrent, found := model.current[key]; found {
				nodeInstance.SetCurrents(entryCurrent)
			}

			for i := 0; i < len(nodeInstance.capasters); i++ {
				if capName, ok := nodeInstance.capasters[i].end.(string); ok {
					if capNode, found := model.nodes[capName]; found {
						nodeInstance.capasters[i].end = capNode
					}
				}
			}

			for i := 0; i < len(nodeInstance.inductance); i++ {
				if indcName, ok := nodeInstance.inductance[i].end.(string); ok {
					if indcNode, found := model.nodes[indcName]; found {
						nodeInstance.inductance[i].end = indcNode
					}
				}
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
				for _, via := range nodeInstance.viases {
					if entryNode, found := model.nodes[via]; found {
						nodeInstance.viasesNodes = append(nodeInstance.viasesNodes, entryNode.connectedNodes...)
						nodeInstance.viasesRes = append(nodeInstance.viasesRes, entryNode.connectedRes...)
					}

					if entryCurrent, found := model.current[via]; found {
						nodeInstance.SetCurrents(entryCurrent)
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
func (model *Model) initElements() {
	bar := prettier.DefaultBar(len(model.nodes), "Initializing the model elements...")

	for _, nodeInstance := range model.nodes {
		nodeInstance.Init()
		bar.Add(1)
	}

	bar.Close()
	fmt.Println()
}

// Restore model.
// Set node values before dc analysis.
func (model *Model) restoreElements() {
	bar := prettier.DefaultBar(len(model.nodes), "Restore the model elements...")

	for _, nodeInstance := range model.nodes {
		nodeInstance.RestorDefaultView()
		bar.Add(1)
	}

	bar.Close()
	fmt.Println()
}

func (model *Model) splitOnDCandTR() ([]string, []string) {
	nodesDC := []string{}
	nodesTR := []string{}

	for _, nodeInstance := range model.nodes {
		if len(nodeInstance.capasters) > 0 {
			nodesTR = append(nodesTR, nodeInstance.name)
		} else {
			nodesDC = append(nodesDC, nodeInstance.name)
		}
	}

	sort.Strings(nodesDC)
	sort.Strings(nodesTR)

	return nodesDC, nodesTR
}

// DC anlysis modeling
func (model *Model) dcModeling() {
	bar := prettier.DefaultBar(model.maxSteps, "Solving the model in dc analysis...")
	total := 0

	model.disableTimeDepentedElements()
	model.collectNodeViases()
	model.prepareElements()
	model.initElements()

	for i := 0; i < model.maxSteps; i++ {
		solvedNodes := 0

		for _, nodeInstance := range model.nodes {
			solvedNodes += nodeInstance.StepDC(model.e, 0)
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

// Transient analysis modeling
func (model *Model) trModeling() {

	// current := NewCurrent("test", 2.18725e-5)
	// current.pulse.max = 1.54613
	// current.pulse.td = 0
	// current.pulse.tr = 1e-10
	// current.pulse.tf = 1e-10
	// current.pulse.pw = 1e-11
	// current.pulse.per = 3e-09

	// for j := 0.0; j < model.analysis.fullTime; j += model.analysis.timeStep {
	// 	fmt.Println(current.PulseValue(j))
	// 	fmt.Println()
	// }

	model.dcModeling()

	fmt.Printf("\n\nInitial values has been calculated.\n")
	fmt.Printf("Starting transient modeling.\n\n")

	model.restoreElements()
	model.collectNodeViases()
	model.prepareElements()
	model.initElements()

	nodesDC, nodesTR := model.splitOnDCandTR()

	// bar := prettier.DefaultBar(int(model.analysis.fullTime/model.analysis.timeStep), "Solving the model in transient analysis...")
	total := 0

	for j := model.analysis.timeStep; j < model.analysis.fullTime; j += model.analysis.timeStep {
		for i := 0; i < model.maxSteps; i++ {
			solvedNodes := 0

			for _, nodeName := range nodesTR {
				if nodeInstance, found := model.nodes[nodeName]; found {
					solvedNodes += nodeInstance.StepTR(model.e, model.analysis.timeStep, j)
				}
			}

			for _, nodeName := range nodesDC {
				if nodeInstance, found := model.nodes[nodeName]; found {
					solvedNodes += nodeInstance.StepTR(model.e, model.analysis.timeStep, j)
				}
			}

			if solvedNodes == len(nodesTR)+len(nodesDC) {
				break
			}
		}

		for _, nodeInstance := range model.nodes {
			nodeInstance.SavePrevV()
		}

		// bar.Add(1)

		total += 1

	}

	for _, nodeInstance := range model.nodes {
		fmt.Println()
		fmt.Println()
		fmt.Println(nodeInstance.name, nodeInstance.prevV)
	}

	// bar.Close()

	fmt.Println()
}

// Modeling system of nodes.
// Modeling last until max steps achived.
func (model *Model) Modeling() {
	switch model.analysis.name {
	case "DC":
		model.dcModeling()
	case "TR":
		model.trModeling()
	}
}
