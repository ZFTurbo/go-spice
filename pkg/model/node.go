package model

import (
	"math"
)

/*
Describe structor of circuits node
*/
type Node struct {
	name               string
	modeling           bool
	current            []*Current
	value              float64
	sumResistors       float64
	connectedNodes     []*Node
	connectedResistors []float64
}

/*
Creates node instanse and return its address.
*/
func NewNode(name string) *Node {
	node := &Node{}
	node.modeling = true
	node.name = name
	node.value = 0.0

	return node
}

/*
Set voltage value
*/
func (n *Node) SetVoltage(value float64) {
	n.value = value
}

/*
Set modeling status
*/
func (n *Node) SetModeling(value bool) {
	n.modeling = value
}

/*
Return curent voltage value of node
*/
func (n *Node) GetVoltage() float64 {
	return n.value
}

/*
Add connected node
*/
func (n *Node) AddNode(node *Node) {
	n.connectedNodes = append(n.connectedNodes, node)
}

/*
* Add connected recistance
 */
func (n *Node) AddResistor(value float64) {
	n.connectedResistors = append(n.connectedResistors, value)
}

/*
* Add connected current source
 */
func (n *Node) AddCurrentSource(value *Current) {
	n.current = append(n.current, value)
}

/*
* Initialize node element.
* Stack connected nodes and res with viases nodes and res, if this node is via.
* Calculates the sum of res and initial value of node voltage.
 */
func (n *Node) Init() {
	if n.modeling {

		sumNodes := 0.0
		sumCurrents := 0.0
		sumResistors := 0.0

		for i := 0; i < len(n.connectedNodes); i++ {
			sumNodes += n.connectedNodes[i].value / n.connectedResistors[i]
		}

		for i := 0; i < len(n.current); i++ {
			sumCurrents += n.current[i].value
		}

		for i := 0; i < len(n.connectedResistors); i++ {
			sumResistors += 1 / n.connectedResistors[i]
		}

		n.value = (sumCurrents - sumNodes) / sumResistors
		n.sumResistors = sumResistors

	}
}

/*
* Make step in modeling for node.
* If previous voltage value and current calculated value differ on less then e then will be return 1.
* Else will be return 0.
 */
func (n *Node) StepDC(e float64) int {
	if n.modeling {
		sumNodes := 0.0
		sumCurrents := 0.0
		lastUpdatedValue := n.value

		for i := 0; i < len(n.connectedNodes); i++ {
			sumNodes += n.connectedNodes[i].value / n.connectedResistors[i]
		}

		for i := 0; i < len(n.current); i++ {
			sumCurrents += n.current[i].value
		}

		n.value = (sumNodes + sumCurrents) / n.sumResistors

		if math.Abs(n.value-lastUpdatedValue) < e {
			return 1
		}

		return 0
	}

	return 1
}
