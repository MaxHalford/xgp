package tree

import (
	"encoding/json"
	"fmt"

	"github.com/MaxHalford/xgp/op"
)

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
func (tr Tree) String() string {
	return CodeDisplay{}.Apply(tr)
}

// Walk recursively applies a function to a Tree and it's branches recursively.
func (tr *Tree) Walk(f func(tr *Tree, depth int) (stop bool)) {
	var apply func(tr *Tree, depth int) bool
	apply = func(tr *Tree, depth int) bool {
		// Apply the function to the Tree and check if the recursion should stop
		if f(tr, depth) {
			return true
		}
		// Apply recursion to each branch
		for _, branch := range tr.branches {
			if apply(branch, depth+1) {
				break
			}
		}
		return false
	}
	apply(tr, 0)
}

// Height returns the height of a tree. The height of a tree is the height of
// its root node. The height of a node is the number of edges on the longest
// path between that node and a leaf.
func (tr Tree) Height() (n int) {
	var f = func(tr *Tree, depth int) (stop bool) {
		if depth > n {
			n = depth
		}
		return
	}
	tr.Walk(f)
	return
}

// count returns the number of nodes that match a specific criteria.
func (tr Tree) count(filter func(*Tree) bool) (n int) {
	var f = func(tr *Tree, depth int) (stop bool) {
		if filter(tr) {
			n++
		}
		return
	}
	tr.Walk(f)
	return n
}

// Size returns the number of Operators in a Tree.
func (tr Tree) Size() int {
	return tr.count(func(*Tree) bool { return true })
}

// NConstants returns the number of Constants in a Tree.
func (tr Tree) NConstants() int {
	var isConstant = func(tr *Tree) bool {
		var _, ok = tr.op.(op.Constant)
		return ok
	}
	return tr.count(isConstant)
}

// Clone a tree by recursively copying it's branches's attributes.
func (tr Tree) Clone() Tree {
	var clone = NewTree(tr.Operator())
	for i, branch := range tr.branches {
		clone.SetBranch(i, branch.Clone())
	}
	return clone
}

// Eval evaluates a Tree on a set of columns.
func (tr Tree) Eval(X [][]float64) (yPred []float64) {
	// If the Tree has no branches then it can be evaluated directly
	if len(tr.branches) == 0 {
		yPred = tr.op.Eval(X)
		return
	}
	// If the Tree has branches then they have to be evaluated first
	var evals = make([][]float64, len(tr.branches))
	for i, branch := range tr.branches {
		evals[i] = branch.Eval(X)
	}
	yPred = tr.op.Eval(evals)
	return
}

// EvalRow is a conveniency method for calling Eval on a single row.
func (tr Tree) EvalRow(x []float64) float64 {
	var X = make([][]float64, len(x))
	for i, xi := range x {
		X[i] = []float64{xi}
	}
	return tr.Eval(X)[0]
}

// EvalRowDebug evaluates a Tree on a single row and displays the successive
// states of the Tree.
func (tr Tree) EvalRowDebug(x []float64, disp Displayer) {
	tr = tr.Clone()
	var (
		height     = tr.Height()
		depthCount = height + 1
	)
	for depthCount >= 0 {
		fmt.Printf("Step %d\n", height-(depthCount-1))
		tr.Walk(func(tr *Tree, depth int) (stop bool) {
			if depthCount == depth {
				tr.op = op.Constant{tr.EvalRow(x)}
				tr.branches = []*Tree{}
			}
			return
		})
		fmt.Println(disp.Apply(tr))
		depthCount--
	}
}

// Simplify a tree by removing unnecessary branches. The algorithm starts at the
// bottom of the tree from left to right. The method returns a boolean to
// indicate if a simplification was performed or not.
func (tr *Tree) Simplify() bool {
	// A tree with no branches can't be simplified
	if len(tr.branches) == 0 {
		return false
	}
	var (
		constBranches = true
		varBranches   = true
	)
	for _, br := range tr.branches {
		// Call the function recursively first so as to start from the bottom
		br.Simplify()
		// Check the type of the branch's operator
		switch br.op.(type) {
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
		tr.op = op.Constant{Value: tr.EvalRow([]float64{0})}
		tr.branches = nil
		return true
	}
	// If the branches are all Variables then a simplification can be made if
	// the mother Operator is of type Sub or of type Div
	if varBranches && tr.NBranches() == 2 {
		// Check if the variables have the same index
		if tr.branches[0].op.(op.Variable).Index == tr.branches[1].op.(op.Variable).Index {
			switch tr.op.(type) {
			case op.Sub:
				tr.op = op.Constant{Value: 0}
				tr.branches = nil
				return true
			case op.Div:
				tr.op = op.Constant{Value: 1}
				tr.branches = nil
				return true
			default:
				return false
			}
		}
	}
	return false
}

// MarshalJSON serializes a Tree into JSON bytes. A serialTree is used as an
// intermediary.
func (tr Tree) MarshalJSON() ([]byte, error) {
	var serial, err = serializeTree(tr)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a Tree. A serialTree is used as an
// intermediary.
func (tr *Tree) UnmarshalJSON(bytes []byte) error {
	var serial serialTree
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedtree, err = parseSerialTree(serial)
	if err != nil {
		return err
	}
	*tr = parsedtree
	return nil
}
