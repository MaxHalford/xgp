package tree

import (
	"math/rand"
)

// A Tree holds an Operator and leaf trees called branches.
type Tree struct {
	Operator Operator
	Branches []*Tree
}

// String representation of a tree.
func (tree *Tree) String() string {
	return CodeDisplay{}.Apply(tree)
}

// rApply recursively applies a function to a Tree and it's branches.
func (tree *Tree) rApply(f func(tree *Tree, depth int) (stop bool)) {
	var apply func(tree *Tree, depth int) bool
	apply = func(tree *Tree, depth int) bool {
		// Apply the function to the Tree and check if the recursion should stop
		if f(tree, depth) {
			return true
		}
		// Apply recursion to each branch
		for _, branch := range tree.Branches {
			if apply(branch, depth+1) {
				break
			}
		}
		return false
	}
	apply(tree, 0)
}

func (tree *Tree) RecApply(f func(tree *Tree)) {
	var w = func(tree *Tree, depth int) (stop bool) { f(tree); return }
	tree.rApply(w)
}

// Height returns the height of a tree. The height of a tree is the height of
// its root node. The height of a node is the number of edges on the longest
// path between that node and a leaf.
func (tree Tree) Height() (n int) {
	var f = func(tree *Tree, depth int) (stop bool) {
		if depth > n {
			n = depth
		}
		return
	}
	tree.rApply(f)
	return
}

// count returns the number of nodes that match a specific criteria.
func (tree Tree) count(filter func(*Tree) bool) (n int) {
	var f = func(tree *Tree, depth int) (stop bool) {
		if filter(tree) {
			n++
		}
		return
	}
	tree.rApply(f)
	return n
}

// NOperators returns the number of Operators in a Tree.
func (tree Tree) NOperators() int {
	return tree.count(func(*Tree) bool { return true })
}

// NConstants returns the number of Constants in a Tree.
func (tree Tree) NConstants() int {
	var isConstant = func(tree *Tree) bool {
		var _, ok = tree.Operator.(Constant)
		return ok
	}
	return tree.count(isConstant)
}

// Clone a tree by recursively copying it's branches's attributes.
func (tree *Tree) Clone() *Tree {
	var branches = make([]*Tree, len(tree.Branches))
	for i, branch := range tree.Branches {
		branches[i] = branch.Clone()
	}
	return &Tree{
		Operator: tree.Operator,
		Branches: branches,
	}
}

// EvaluateRow blabla
func (tree Tree) EvaluateRow(x []float64) float64 {
	// Either the tree is a leaf tree
	if len(tree.Branches) == 0 {
		return tree.Operator.ApplyRow(x)
	}
	// Either the tree has branches trees
	var evals = make([]float64, len(tree.Branches))
	for i, branch := range tree.Branches {
		evals[i] = branch.EvaluateRow(x)
	}
	return tree.Operator.ApplyRow(evals)
}

// EvaluateCols blabla
func (tree *Tree) EvaluateCols(X [][]float64, cache *Cache) ([]float64, error) {
	// Check the cache
	if cache != nil {
		var yPred = cache.Get(tree.String())
		if yPred != nil {
			return yPred, nil
		}
	}
	// Simplify the tree to remove unnecessary evaluation parts
	tree.simplify()
	// The Tree has no branches
	if len(tree.Branches) == 0 {
		return tree.Operator.ApplyCols(X), nil
	}
	// The Tree has branches
	var evals = make([][]float64, len(tree.Branches))
	for i, branch := range tree.Branches {
		var eval, err = branch.EvaluateCols(X, cache)
		if err != nil {
			return nil, err
		}
		evals[i] = eval
	}
	var yPred = tree.Operator.ApplyCols(evals)
	// Add the results to the cache
	if cache != nil {
		cache.Set(tree.String(), yPred)
	}
	return yPred, nil
}

// setOperator replaces the Operator of a tree.
func (tree *Tree) setOperator(op Operator, rng *rand.Rand) {
	tree.Operator = op
}

// Simplify a tree by removing unnecessary branches. The algorithm starts at the
// bottom of the tree from left to right. The method returns a boolean to
// indicate if a simplification was performed or not.
func (tree *Tree) simplify() bool {
	// A tree with no branches can't be simplified
	if len(tree.Branches) == 0 {
		return false
	}
	var constBranches = true
	for i, branch := range tree.Branches {
		// Call the function recursively first so as to start from the bottom
		tree.Branches[i].simplify()
		// Check if the tree has a branch that isn't a Constant
		if _, ok := branch.Operator.(Constant); !ok {
			constBranches = false
		}
	}
	// Stop if the tree has no branches with Variable Operators
	if !constBranches {
		return false
	}
	// Replace the tree's Operator with a Constant.
	tree.Operator = Constant{Value: tree.EvaluateRow([]float64{})}
	tree.Branches = nil
	return true
}
