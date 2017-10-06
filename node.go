package xgp

import (
	"errors"
	"math/rand"

	"github.com/MaxHalford/xgp/tree"
	"gonum.org/v1/gonum/floats"
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

// NConstants returns the number of Constants in a Node tree.
func (node Node) NConstants() int {
	var (
		n int
		f = func(node *Node) {
			if _, ok := node.Operator.(Constant); ok {
				n++
			}
		}
	)
	node.RecApply(f)
	return n
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
func (node Node) evaluateRow(x []float64) float64 {
	// Either the Node is a leaf Node
	if len(node.Children) == 0 {
		return node.Operator.Apply(x)
	}
	// Either the Node has children Nodes
	var childEvals = make([]float64, len(node.Children))
	for i, child := range node.Children {
		childEvals[i] = child.evaluateRow(x)
	}
	return node.Operator.Apply(childEvals)
}

func (node *Node) evaluateXT(XT [][]float64) (y []float64, err error) {
	// Simplify the Node to remove unnecessary evaluation parts
	node.Simplify()
	// Either the Node is a leaf Node
	var childEvals = make([][]float64, len(node.Children))
	if node.IsTerminal() {
		y = node.Operator.ApplyXT(XT)
	} else {
		// Either the Node has children Nodes
		for i, child := range node.Children {
			childEvals[i], err = child.evaluateXT(XT)
			if err != nil {
				return nil, err
			}
		}
		y = node.Operator.ApplyXT(childEvals)
	}
	// Store the result
	if floats.HasNaN(y) {
		return nil, errors.New("Slice contains NaNs")
	}
	return
}

// setOperator replaces the Operator of a Node.
func (node *Node) setOperator(op Operator, rng *rand.Rand) {
	node.Operator = op
}

// String representation of a Node.
func (node *Node) String() string {
	var displayer = tree.EquationDisplay{}
	return displayer.Apply(node)
}

// Simplify a Node by removing unnecessary children. The algorithm starts at the
// bottom of the tree from left to right.
func (node *Node) Simplify() {
	// A Node with no children can't be simplified
	if node.NBranches() == 0 {
		return
	}
	var constChildren = true
	for i, child := range node.Children {
		// Call the function recursively first so as to start from the bottom
		node.Children[i].Simplify()
		// Check if the Node has a child that isn't a Constant
		if _, ok := child.Operator.(Constant); !ok {
			constChildren = false
		}
	}
	// Stop if the Node has no children with Variable Operators
	if !constChildren {
		return
	}
	// Replace the Node's Operator with a Constant.
	node.Operator = Constant{Value: node.evaluateRow([]float64{})}
	node.Children = nil
}

// IsTerminal indicates if a Node is a terminal Node or not.
func (node Node) IsTerminal() bool { return node.NBranches() == 0 }

// NBranches is required to implement tree.Tree.
func (node *Node) NBranches() int { return len(node.Children) }

// GetBranch is required to implement tree.Tree.
func (node *Node) GetBranch(i int) tree.Tree { return node.Children[i] }

// Swap is required to implement tree.Tree.
func (node *Node) Swap(otherTree tree.Tree) { *node, *otherTree.(*Node) = *otherTree.(*Node), *node }

// ToString is required to implement tree.Tree.
func (node *Node) ToString() string { return node.Operator.String() }
