package model

import (
	"math"
)

// Describe structor of circuits node
type Node struct {
	name           string
	modeling       bool
	i              []*Current
	v              float64
	dV             float64
	prevV          []float64
	sumRes         float64
	connectedNodes []interface{}
	connectedRes   []float64
	viases         []string
	viasesRes      []float64
	viasesNodes    []interface{}
	inductance     []*Inductance
	capasters      []*Capaster
}

// Creates node instanse and return its address.
func NewNode(name string, end_start interface{}, resistance float64) *Node {
	node := &Node{}
	node.modeling = true
	node.name = name
	node.connectedNodes = append(node.connectedNodes, end_start)
	node.connectedRes = append(node.connectedRes, resistance)

	return node
}

func NewNodeRaw(name string) *Node {
	node := &Node{}
	node.name = name
	node.modeling = true

	return node
}

// Creates via instanse and return its address.
func NewVia(name string, via_name string) *Node {
	via := &Node{}
	via.modeling = true
	via.name = name
	via.viases = append(via.viases, via_name)

	return via
}

// Return curent voltage value of node
func (n *Node) GetVoltage() float64 {
	return n.v
}

// Add connected node
func (n *Node) AddNode(name string) {
	n.connectedNodes = append(n.connectedNodes, name)
}

// Add connected recistance
func (n *Node) AddRes(val float64) {
	n.connectedRes = append(n.connectedRes, val)
}

// Add connected via, if not exists in nodes viases
// Return false if exists else return true
func (n *Node) AddVia(name string) bool {
	for _, via := range n.viases {
		if via == name {
			return false
		}
	}

	n.viases = append(n.viases, name)

	return true
}

// Add induction
func (n *Node) AddInductance(val *Inductance) {
	n.inductance = append(n.inductance, val)
}

// Add capaster
func (n *Node) AddCapaster(val *Capaster) {
	n.capasters = append(n.capasters, val)
}

// Add current source
func (n *Node) SetCurrents(c []*Current) {
	n.i = c
}

// Set voltage value
func (n *Node) SetVoltage(val float64) {
	n.v = val
}

// Set modeling status
func (n *Node) SetModeling(val bool) {
	n.modeling = val
}

// Save prev value in tran analysis
func (n *Node) SavePrevV() {
	n.prevV = append(n.prevV, n.v)
}

// Calculates the sum of the elements of an array of the form 1/x.
// Where x is element of the array.
func sumReverse(array []float64) float64 {
	sum := 0.0

	for i := 0; i < len(array); i++ {
		sum += 1 / array[i]
	}

	return sum
}

// Calculates the sum of the elements of two arrays of the form y/x.
// Where x is element of the first array and y is element of second array.
func sumZip(y []interface{}, x []float64) float64 {
	sum := 0.0

	for i, y := range y {
		if entryY, ok := y.(float64); ok {
			sum += entryY / x[i]
		}
	}

	return sum
}

// Initialize node element.
// Stack connected nodes and res with viases nodes and res, if this node is via.
// Calculates the sum of res and initial value of node voltage.
func (n *Node) Init() {
	if n.modeling {
		if n.v == 0 {
			sumI := 0.0

			for i := 0; i < len(n.i); i++ {
				sumI += n.i[i].val
			}

			n.v = sumI - sumZip(n.connectedNodes, n.connectedRes)
		}

		n.connectedNodes = append(n.connectedNodes, n.viasesNodes...)
		n.connectedRes = append(n.connectedRes, n.viasesRes...)
		n.sumRes = sumReverse(n.connectedRes)
		n.prevV = []float64{n.v}
	}

}

// Restor all node dependites, expect node voltage value.
func (n *Node) RestorDefaultView() {
	oldConnectedNodes := []interface{}{}
	oldConnectedRes := []float64{}
	oldViases := []string{}

	for i, node := range n.connectedNodes {
		if val1, ok := node.(*Node); ok {
			viaConnection := false

			for _, nodeVia := range n.viasesNodes {
				if val2, ok := nodeVia.(*Node); ok && val2.name == val1.name {
					viaConnection = true
				}
			}

			if !viaConnection {
				oldConnectedNodes = append(oldConnectedNodes, val1.name)
				oldConnectedRes = append(oldConnectedRes, n.connectedRes[i])
			}
		}

		if val1, ok := node.(float64); ok {
			oldConnectedNodes = append(oldConnectedNodes, val1)
			oldConnectedRes = append(oldConnectedRes, n.connectedRes[i])
		}
	}

	for _, via := range n.viases {
		viaInIndc := false

		for _, indc := range n.inductance {
			if val, ok := indc.end.(*Node); ok && val.name == via {
				viaInIndc = true
				break
			}
		}

		if !viaInIndc {
			oldViases = append(oldViases, via)
		}
	}

	n.modeling = true
	n.i = []*Current{}
	n.connectedNodes = oldConnectedNodes
	n.connectedRes = oldConnectedRes
	n.viases = oldViases
	n.viasesNodes = []interface{}{}
	n.viasesRes = []float64{}

}

// Make step in modeling for node.
// If previous voltage value and current calculated value differ on less then e then will be return 1.
// Else will be return 0.
func (n *Node) StepDC(e float64, t float64) int {
	if n.modeling {
		_v := n.v
		sum := 0.0
		sumI := 0.0

		for i := 0; i < len(n.connectedNodes); i++ {
			if entryNode, ok := n.connectedNodes[i].(float64); ok {
				sum += entryNode / n.connectedRes[i]
			}
			if entryNode, ok := n.connectedNodes[i].(*Node); ok {
				sum += entryNode.v / n.connectedRes[i]
			}
		}

		for i := 0; i < len(n.i); i++ {
			val1, _ := n.i[i].PulseValue(t)
			sumI += val1
		}

		n.v = 1.75*(sum-sumI)/n.sumRes + (1-1.75)*_v

		if math.Abs(n.v-_v) < e {
			return 1
		}

		return 0
	}

	return 1
}

func (n *Node) StepTR(e float64, h float64, t float64) int {
	if n.modeling {
		
	}

	return 1
}
