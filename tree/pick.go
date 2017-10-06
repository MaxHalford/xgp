package tree

import (
	"math/rand"
	"sort"

	"gonum.org/v1/gonum/floats"
)

func PickSubTree(tree Tree, weight func(tree Tree) float64, minDepth, maxDepth int, rng *rand.Rand) (Tree, int) {
	// Assign weight to each Tree and calculate the total weight
	var (
		weights      []float64
		totalWeight  float64
		assignWeight = func(tree Tree, depth int) (stop bool) {
			var w float64
			if depth < minDepth || (depth > maxDepth && maxDepth >= 0) {
				w = 0
			} else {
				w = weight(tree) + 1
			}
			weights = append(weights, w)
			totalWeight += w
			return
		}
	)
	rApply(tree, assignWeight)
	// Cumulatively sum the weights
	var cumsum = make([]float64, len(weights))
	floats.CumSum(cumsum, weights)
	// Sample a random number in [0, cumsum[-1])
	var r = rng.Float64() * cumsum[len(cumsum)-1]
	// Find i where cumsum[i-1] < r < cumsum[i]
	var (
		pos      = sort.SearchFloat64s(cumsum, r)
		posDepth int
	)
	// Extract the sub-tree at position i
	var (
		subTree     Tree
		i           int
		findSubTree = func(tree Tree, depth int) (stop bool) {
			if i < pos {
				i++
				return
			}
			subTree = tree
			posDepth = depth
			return true
		}
	)
	rApply(tree, findSubTree)
	return subTree, posDepth
}
