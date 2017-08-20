package xgp

import (
	"math/rand"

	"github.com/MaxHalford/xgp/tree"
)

// A Node holds an Operator and leaf Nodes called children. A Node is said to be
// terminal if it has no children. The Operator's arity and the number of
// children should always be the same.
type Node struct {
	Operator Operator
	Children []*Node
}

// RecApply recursively applies a function to a Node and it's children.
func (node *Node) RecApply(f func(*Node)) {
	f(node)
	for _, child := range node.Children {
		child.RecApply(f)
	}
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

// Simplify a Node by removing unnecessary children. The algorithm starts at the
// bottom of the tree from left to right.
func (node *Node) Simplify() {
	var varChildren bool
	// Call the function recursively first so as to start from the bottom
	for i, child := range node.Children {
		node.Children[i].Simplify()
		// Check if the Node has children that contain Variable Operators
		if _, ok := child.Operator.(Variable); ok {
			varChildren = true
		}
	}
	// Stop if the Node has no children or one of them has a Variable Operator
	if varChildren || node.NBranches() == 0 {
		return
	}
	// Replace the Node's Operator with a Constant.
	node.Operator = Constant{Value: node.evaluate([]float64{})}
	node.Children = nil
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

// ToString method is required to implement the Tree interface from the tree
// package.
func (node *Node) ToString() string {
	return node.Operator.String()
}
