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
	Weighting Weighting
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
				w = wp.Weighting.apply(tree.Operator)
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
