package xgp

import "math/rand"

type newNode func(leaf bool, rng *rand.Rand) *Node

// A NodeInitializer generates a random Node.
type NodeInitializer interface {
	Apply(newNode newNode, rng *rand.Rand) *Node
}

// FullNodeInitializer generates Nodes where all the leaves are the same depth.
type FullNodeInitializer struct {
	Height int
}

func (init FullNodeInitializer) Apply(newNode newNode, rng *rand.Rand) *Node {
	if init.Height == 0 {
		return newNode(true, rng)
	}
	var node = newNode(false, rng)
	for i := range node.Children {
		node.Children[i] = FullNodeInitializer{Height: init.Height - 1}.Apply(newNode, rng)
	}
	return node
}

// GrowNodeInitializer generates Nodes where all the leaves have at most a
// certain depth.
type GrowNodeInitializer struct {
	MaxHeight int
	PLeaf     float64 // Probability of producing a leaf
}

func (init GrowNodeInitializer) Apply(newNode newNode, rng *rand.Rand) *Node {
	if init.MaxHeight == 0 || rng.Float64() < init.PLeaf {
		return newNode(true, rng)
	}
	var node = newNode(false, rng)
	for i := range node.Children {
		node.Children[i] = GrowNodeInitializer{
			MaxHeight: init.MaxHeight - 1,
			PLeaf:     init.PLeaf,
		}.Apply(newNode, rng)
	}
	return node
}
