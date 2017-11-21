package tree

import (
	"math/rand"
	"sort"

	"gonum.org/v1/gonum/floats"
)

// A Picker picks a sub-Tree from a Tree.
type Picker interface {
	Apply(tree *Tree, minDepth, maxDepth int, rng *rand.Rand) (subTree *Tree, subTreeDepth int)
}

// WeightPicker picks a sub-Tree at random by weighting each sub-tree
// according to it's Operator's type.
type WeightedPicker struct {
	PConstant float64
	PVariable float64
	PFunction float64
}

// weight is just a utility method to get the weight that matches an Operator's
// type.
func (wp WeightedPicker) weight(op Operator) float64 {
	switch op.(type) {
	case Constant:
		return wp.PConstant
	case Variable:
		return wp.PVariable
	default:
		return wp.PFunction
	}
}

// Apply WeightedPicker.
func (wp WeightedPicker) Apply(tree *Tree, minDepth, maxDepth int, rng *rand.Rand) (*Tree, int) {
	// Assign weight to each Tree and calculate the total weight
	var (
		trees        []*Tree
		depths       []int
		weights      []float64
		totalWeight  float64
		assignWeight = func(tree *Tree, depth int) (stop bool) {
			var w float64
			if depth < minDepth || (depth > maxDepth && maxDepth >= 0) {
				w = 0
			} else {
				w = wp.weight(tree.Operator)
			}
			weights = append(weights, w)
			trees = append(trees, tree)
			depths = append(depths, depth)
			totalWeight += w
			return
		}
	)
	tree.rApply(assignWeight)
	// Calculate the cumulative sum of the weights
	var cumSum = make([]float64, len(weights))
	floats.CumSum(cumSum, weights)
	// Sample a random number in [0, cumSum[-1])
	var r = rng.Float64() * cumSum[len(cumSum)-1]
	// Find i where cumSum[i-1] < r < cumSum[i]
	var i = sort.SearchFloat64s(cumSum, r)
	return trees[i], depths[i]
}
