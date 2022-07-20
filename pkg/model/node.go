package model

import (
	"fmt"
	"math"
)

// Describe structor of circuits node
type Node struct {
	name           string
	i              []*Current
	v              float64
	prevV          float64
	sumRes         float64
	connectedNodes []interface{}
	connectedRes   []float64
	viases         []string
	viasesRes      []float64
	viasesNodes    []interface{}
	inductance     []float64
	capasters      []float64
}

// Creates node instanse and return its address.
func NewNode(name string, end_start interface{}, resistance float64) *Node {
	node := &Node{}
	node.name = name
	node.connectedNodes = append(node.connectedNodes, end_start)
	node.connectedRes = append(node.connectedRes, resistance)
	return node
}

// Creates via instanse and return its address.
func NewVia(name string, via_name string) *Node {
	via := &Node{}
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

// Add connected via
func (n *Node) AddVia(name string) {
	n.viases = append(n.viases, name)
}

// Add induction
func (n *Node) AddInductance(val float64) {
	n.inductance = append(n.inductance, val)
}

// Add capaster
func (n *Node) AddCupaster(val float64) {
	n.capasters = append(n.capasters, val)
}

// Add current source
func (n *Node) SetCurrents(c []*Current) {
	n.i = c
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
	sumI := 0.0

	for i := 0; i < len(n.i); i++ {
		sumI += n.i[i].val
	}

	n.prevV = 0
	n.connectedNodes = append(n.connectedNodes, n.viasesNodes...)
	n.connectedRes = append(n.connectedRes, n.viasesRes...)
	n.sumRes = sumReverse(n.connectedRes)
	n.v = sumI - sumZip(n.connectedNodes, n.connectedRes)

	n.viases = nil
	n.viasesRes = nil
	n.viasesNodes = nil
}

func (n *Node) SavePrevTimeStep() {
	n.prevV = n.v
}

// Make step in modeling for node.
// If previous voltage value and current calculated value differ on less then e then will be return 1.
// Else will be return 0.
// Note that node and res count must be equivalent.
func (n *Node) Step(e float64, h float64, t float64) int {
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
		sumI += n.i[i].PulseValue(t)
	}

	n.v = (sum - sumI) / n.sumRes

	if len(n.capasters) != 0 {
		capSum := 0.0

		for i := 0; i < len(n.capasters); i++ {
			capSum += 1 / n.capasters[i]
		}

		n.v = n.prevV + h*(-sum+sumI+n.sumRes*n.prevV)*capSum

		if n.name == "n0_2679_17913" {
			fmt.Println(n.capasters)
		}
	}

	n.v = 1.75*n.v + (1-1.75)*_v

	if math.Abs(n.v-_v) < e {
		return 1
	}

	return 0
}
