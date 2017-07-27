package xgp

import (
	"math/rand"

	"github.com/MaxHalford/xgp/tree"
)

type Node struct {
	Operator Operator
	Children []*Node
}

// Clone a Node by recursively copying it's children's attributes.
func (node *Node) clone() *Node {
	var children = make([]*Node, len(node.Children))
	for i, child := range node.Children {
		children[i] = child.clone()
	}
	return &Node{
		Operator: node.Operator,
		Children: children,
	}
}

// Evaluate a Node by evaluating it's children recursively and running the
// children's output through the Node's Operator.
func (node Node) evaluate(X []float64) float64 {
	// Either the Node is a leaf Node
	if len(node.Children) == 0 {
		return node.Operator.Apply(X)
	}
	// Either the Node is an internal Node
	var childEvals = make([]float64, len(node.Children))
	for i, child := range node.Children {
		childEvals[i] = child.evaluate(X)
	}
	return node.Operator.Apply(childEvals)
}

// setOperator replaces the Operator of a Node.
func (node *Node) setOperator(op Operator, rng *rand.Rand) {
	node.Operator = op
}

// String representation of a Node.
func (node *Node) String() string {
	var displayer = tree.DirDisplay{TabSize: 4}
	return displayer.Apply(node)
}

// Implementation of the Tree interface from the tree package

// NBranches method is required to implement the Tree interface from the tree
// package.
func (node *Node) NBranches() int {
	return len(node.Children)
}

// GetBranch method is required to implement to Tree interface from the tree
// package.
func (node *Node) GetBranch(i int) tree.Tree {
	return node.Children[i]
}

// Swap method is required to implement to Tree interface from the tree package.
func (node *Node) Swap(otherTree tree.Tree) {
	*node, *otherTree.(*Node) = *otherTree.(*Node), *node
}

// String method is required to implement the Tree interface from the tree
// package.
func (node *Node) ToString() string {
	return node.Operator.String()
}
