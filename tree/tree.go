package tree

import "github.com/MaxHalford/koza/op"

// A Tree holds an Operator and branches.
type Tree struct {
	op       op.Operator
	branches []*Tree
}

// Operator returns a Tree's Operator.
func (tr Tree) Operator() op.Operator {
	return tr.op
}

// SetBranch replaces the ith branch with a given Tree.
func (tr *Tree) SetBranch(i int, br Tree) {
	tr.branches[i] = &br
}

// SetOperator replaces a Tree's Operator.
func (tr *Tree) SetOperator(op op.Operator) {
	tr.op = op
}

// NBranches returns the number of branches of a Tree.
func (tr Tree) NBranches() int {
	return len(tr.branches)
}

// NewTree returns a Tree with a given Operator.
func NewTree(op op.Operator) Tree {
	return Tree{
		op:       op,
		branches: make([]*Tree, op.Arity()),
	}
}

// String representation of a tree.
func (tree Tree) String() string {
	return CodeDisplay{}.Apply(tree)
}

// Walk recursively applies a function to a Tree and it's branches recursively.
func (tree *Tree) Walk(f func(tree *Tree, depth int) (stop bool)) {
	var apply func(tree *Tree, depth int) bool
	apply = func(tree *Tree, depth int) bool {
		// Apply the function to the Tree and check if the recursion should stop
		if f(tree, depth) {
			return true
		}
		// Apply recursion to each branch
		for _, branch := range tree.branches {
			if apply(branch, depth+1) {
				break
			}
		}
		return false
	}
	apply(tree, 0)
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
	tree.Walk(f)
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
	tree.Walk(f)
	return n
}

// Size returns the number of Operators in a Tree.
func (tree Tree) Size() int {
	return tree.count(func(*Tree) bool { return true })
}

// NConstants returns the number of Constants in a Tree.
func (tree Tree) NConstants() int {
	var isConstant = func(tree *Tree) bool {
		var _, ok = tree.op.(op.Constant)
		return ok
	}
	return tree.count(isConstant)
}

// Clone a tree by recursively copying it's branches's attributes.
func (tr Tree) Clone() Tree {
	var clone = NewTree(tr.Operator())
	for i, branch := range tr.branches {
		clone.SetBranch(i, branch.Clone())
	}
	return clone
}

// evaluateRow blabla
func (tree Tree) EvaluateRow(x []float64) float64 {
	// Either the tree is a leaf tree
	if len(tree.branches) == 0 {
		return tree.op.ApplyRow(x)
	}
	// Either the tree has branches that have to be evaluated first
	var evals = make([]float64, len(tree.branches))
	for i, branch := range tree.branches {
		evals[i] = branch.EvaluateRow(x)
	}
	return tree.op.ApplyRow(evals)
}

// EvaluateCols blabla
func (tree *Tree) EvaluateCols(X [][]float64) (yPred []float64, err error) {
	// Simplify the tree to remove unnecessary evaluation parts
	tree.simplify()

	// If the Tree has no branches then it can be evaluated directly
	if len(tree.branches) == 0 {
		yPred = tree.op.ApplyCols(X)
		return
	}

	// If the Tree has branches then they have to be evaluated first
	var evals = make([][]float64, len(tree.branches))
	for i, branch := range tree.branches {
		var eval, err = branch.EvaluateCols(X)
		if err != nil {
			return nil, err
		}
		evals[i] = eval
	}
	yPred = tree.op.ApplyCols(evals)
	return
}

// Simplify a tree by removing unnecessary branches. The algorithm starts at the
// bottom of the tree from left to right. The method returns a boolean to
// indicate if a simplification was performed or not.
func (tree *Tree) simplify() bool {
	// A tree with no branches can't be simplified
	if len(tree.branches) == 0 {
		return false
	}
	var (
		constBranches = true
		varBranches   = true
	)
	for _, branch := range tree.branches {
		// Call the function recursively first so as to start from the bottom
		branch.simplify()
		// Check the type of the branch's operator
		switch branch.op.(type) {
		case op.Constant:
			varBranches = false
		case op.Variable:
			constBranches = false
		default:
			varBranches = false
			constBranches = false
		}
	}
	// If the branches are all Constants then a simplification can be made
	if constBranches {
		tree.op = op.Constant{Value: tree.EvaluateRow([]float64{})}
		tree.branches = nil
		return true
	}
	// If the branches are all Variables then a simplification can be made if
	// the mother Operator is of type Difference
	if varBranches && tree.NBranches() == 2 {
		// Check if the variables have the same index
		if tree.branches[0].op.(op.Variable).Index == tree.branches[1].op.(op.Variable).Index {
			switch tree.op.(type) {
			case op.Difference:
				tree.op = op.Constant{Value: 0}
				tree.branches = nil
				return true
			case op.Division:
				tree.op = op.Constant{Value: 1}
				tree.branches = nil
				return true
			default:
				return false
			}
		}
	}
	return false
}
