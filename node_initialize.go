package xgp

import "math/rand"

type newOperator func(terminal bool, rng *rand.Rand) Operator

// A NodeInitializer generates a random Node.
type NodeInitializer interface {
	Apply(newOperator newOperator, rng *rand.Rand) *Node
}

// FullNodeInitializer generates Nodes where all the leaves are the same depth.
type FullNodeInitializer struct {
	Height int
}

// Apply FullNodeInitializer.
func (init FullNodeInitializer) Apply(newOperator newOperator, rng *rand.Rand) *Node {
	var op Operator
	if init.Height == 0 {
		op = newOperator(true, rng)
	} else {
		op = newOperator(false, rng)
	}
	var node = &Node{
		Operator: op,
		Children: make([]*Node, op.Arity()),
	}
	for i := range node.Children {
		node.Children[i] = FullNodeInitializer{Height: init.Height - 1}.Apply(newOperator, rng)
	}
	return node
}

// GrowNodeInitializer generates Nodes where all the leaves have at most a
// certain depth.
type GrowNodeInitializer struct {
	MaxHeight int
	PLeaf     float64 // Probability of producing a leaf
}

// Apply GrowNodeInitializer.
func (init GrowNodeInitializer) Apply(newOperator newOperator, rng *rand.Rand) *Node {
	var op Operator
	if init.MaxHeight == 0 || rng.Float64() < init.PLeaf {
		op = newOperator(true, rng)
	} else {
		op = newOperator(false, rng)
	}
	var node = &Node{
		Operator: op,
		Children: make([]*Node, op.Arity()),
	}
	for i := range node.Children {
		node.Children[i] = GrowNodeInitializer{
			MaxHeight: init.MaxHeight - 1,
			PLeaf:     init.PLeaf,
		}.Apply(newOperator, rng)
	}
	return node
}
