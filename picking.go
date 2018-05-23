package xgp

import (
	"math/rand"
	"sort"

	"github.com/MaxHalford/xgp/op"
	"github.com/MaxHalford/xgp/tree"
	"gonum.org/v1/gonum/floats"
)

// A Weighting is a convenience structure for assigning weights to Operators
// for selection purposes.
type Weighting struct {
	PConstant float64
	PVariable float64
	PFunction float64
}

func (w Weighting) apply(operator op.Operator) float64 {
	switch operator.(type) {
	case op.Constant:
		return w.PConstant
	case op.Variable:
		return w.PVariable
	default:
		return w.PFunction
	}
}

// A Picker picks a subtree from a Tree. The subtree can be forced to have
// a height in range [minHeight, maxHeight].
type Picker interface {
	Apply(tree *tree.Tree, minHeight, maxHeight int, rng *rand.Rand) *tree.Tree
}

// WeightedPicker picks a subtree at random by weighting each subtree
// according to it's Operator's type.
type WeightedPicker struct {
	Weighting Weighting
}

// Apply WeightedPicker.
func (wp WeightedPicker) Apply(tr *tree.Tree, minHeight, maxHeight int, rng *rand.Rand) *tree.Tree {
	// Assign weight to each Tree and calculate the total weight
	var (
		trees        []*tree.Tree
		weights      []float64
		totalWeight  float64
		assignWeight = func(tr *tree.Tree, depth int) (stop bool) {
			var (
				w float64
				h = tr.Height()
			)
			if h < minHeight || h > maxHeight {
				w = 0
			} else {
				w = wp.Weighting.apply(tr.Op)
			}
			weights = append(weights, w)
			trees = append(trees, tr)
			totalWeight += w
			return
		}
	)
	tr.Walk(assignWeight)

	// Calculate the cumulative sum of the weights
	var cumSum = make([]float64, len(weights))
	floats.CumSum(cumSum, weights)

	// Sample a random number in [0, cumSum[-1])
	var r = rng.Float64() * cumSum[len(cumSum)-1]

	// Find i where cumSum[i-1] < r < cumSum[i]
	var i = sort.SearchFloat64s(cumSum, r)

	return trees[i]
}
