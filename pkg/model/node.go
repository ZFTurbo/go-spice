package model

import (
	"math"
)

// Describe structor of circuits node
type Node struct {
	Modeling       bool
	Name           string
	I              float64
	V              float64
	PrevV          float64
	SumRes         float64
	ConnectedNodes []interface{}
	ConnectedRes   []float64
	Viases         []string
	ViasesRes      []float64
	ViasesNodes    []interface{}
}

// Creates node instanse and return its address.
func NewNode(name string, end_start interface{}, resistance float64) *Node {
	node := &Node{}
	node.ConnectedNodes = append(node.ConnectedNodes, end_start)
	node.ConnectedRes = append(node.ConnectedRes, resistance)
	return node
}

// Creates via instanse and return its address.
func NewVia(name string, via_name string) *Node {
	via := &Node{}
	via.Viases = append(via.Viases, via_name)
	return via
}

// Calculates the sum of the elements of an array of the form 1/x.
// Where x is element of the array.
func SumReverse(array []float64) float64 {
	sum := 0.0

	for i := 0; i < len(array); i++ {
		sum += 1 / array[i]
	}

	return sum
}

// Calculates the sum of the elements of two arrays of the form y/x.
// Where x is element of the first array and y is element of second array.
func SumZip(y []interface{}, x []float64) float64 {
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
	if n.Modeling {
		n.ConnectedNodes = append(n.ConnectedNodes, n.ViasesNodes...)
		n.ConnectedRes = append(n.ConnectedRes, n.ViasesRes...)
		n.SumRes = SumReverse(n.ConnectedRes)
		n.V = n.I - SumZip(n.ConnectedNodes, n.ConnectedRes)
		n.Viases = nil
		n.ViasesRes = nil
		n.ViasesNodes = nil
	}
}

// Make step in modeling for node.
// If previous voltage value and current calculated value differ on less then e then will be return 1.
// Else will be return 0.
// Note that node and res count must be equivalent.
func (n *Node) Step(e float64) int {
	if n.Modeling {
		n.PrevV = n.V
		sum := 0.0

		for i := 0; i < len(n.ConnectedNodes); i++ {
			if entryNode, ok := n.ConnectedNodes[i].(float64); ok {
				sum += entryNode / n.ConnectedRes[i]
			}
			if entryNode, ok := n.ConnectedNodes[i].(*Node); ok {
				sum += entryNode.V / n.ConnectedRes[i]
			}
		}

		n.V = sum/n.SumRes - n.I/n.SumRes

		if math.Abs(n.V-n.PrevV) < e {
			return 1
		}

		return 0
	} else {
		return 1
	}
}
