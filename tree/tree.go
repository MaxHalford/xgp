package tree

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MaxHalford/xgp/op"
)

// A Tree holds an Operator and branches.
type Tree struct {
	Op       op.Operator
	Branches []*Tree
}

// NewTree returns a Tree with a given Operator.
func NewTree(op op.Operator) Tree {
	return Tree{
		Op:       op,
		Branches: make([]*Tree, op.Arity()),
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
		for _, branch := range tr.Branches {
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
		var _, ok = tr.Op.(op.Constant)
		return ok
	}
	return tr.count(isConstant)
}

// Clone a tree by recursively copying it's branches's attributes.
func (tr Tree) Clone() Tree {
	var clone = NewTree(tr.Op)
	for i, br := range tr.Branches {
		brc := br.Clone()
		clone.Branches[i] = &brc
	}
	return clone
}

// Eval evaluates a Tree on a set of columns.
func (tr Tree) Eval(X [][]float64) (yPred []float64) {
	// If the Tree has no branches then it can be evaluated directly
	if len(tr.Branches) == 0 {
		yPred = tr.Op.Eval(X)
		return
	}
	// If the Tree has branches then they have to be evaluated first
	var evals = make([][]float64, len(tr.Branches))
	for i, branch := range tr.Branches {
		evals[i] = branch.Eval(X)
	}
	yPred = tr.Op.Eval(evals)
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
				tr.Op = op.Constant{tr.EvalRow(x)}
				tr.Branches = []*Tree{}
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
	if len(tr.Branches) == 0 {
		return false
	}
	var (
		constBranches = true
		varBranches   = true
	)
	for _, br := range tr.Branches {
		// Call the function recursively first so as to start from the bottom
		br.Simplify()
		// Check the type of the branch's operator
		switch br.Op.(type) {
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
		tr.Op = op.Constant{Value: tr.EvalRow([]float64{0})}
		tr.Branches = nil
		return true
	}
	// If the branches are all Variables then a simplification can be made if
	// the mother Operator is of type Sub or of type Div
	if varBranches && len(tr.Branches) == 2 {
		// Check if the variables have the same index
		if tr.Branches[0].Op.(op.Variable).Index == tr.Branches[1].Op.(op.Variable).Index {
			switch tr.Op.(type) {
			case op.Sub:
				tr.Op = op.Constant{Value: 0}
				tr.Branches = nil
				return true
			case op.Div:
				tr.Op = op.Constant{Value: 1}
				tr.Branches = nil
				return true
			default:
				return false
			}
		}
	}
	return false
}

// MarshalJSON serializes a Tree into JSON bytes.
func (tr Tree) MarshalJSON() ([]byte, error) {
	var serial = struct {
		OpType   string  `json:"op_type"`
		OpValue  string  `json:"op_value"`
		Branches []*Tree `json:"branches"`
	}{
		Branches: tr.Branches,
	}

	switch tr.Op.(type) {
	case op.Constant:
		serial.OpType = "constant"
		serial.OpValue = strconv.FormatFloat(tr.Op.(op.Constant).Value, 'f', -1, 64)
	case op.Variable:
		serial.OpType = "variable"
		serial.OpValue = strconv.Itoa(tr.Op.(op.Variable).Index)
	default:
		serial.OpType = "function"
		serial.OpValue = tr.Op.String()
	}

	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a Tree.
func (tr *Tree) UnmarshalJSON(bytes []byte) error {
	var serial = &struct {
		OpType   string  `json:"op_type"`
		OpValue  string  `json:"op_value"`
		Branches []*Tree `json:"branches"`
	}{}

	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}

	switch serial.OpType {
	case "constant":
		var val, err = strconv.ParseFloat(serial.OpValue, 64)
		if err != nil {
			return err
		}
		tr.Op = op.Constant{val}
	case "variable":
		var idx, err = strconv.Atoi(serial.OpValue)
		if err != nil {
			return err
		}
		tr.Op = op.Variable{idx}
	default:
		var function, err = op.ParseFunc(serial.OpValue)
		if err != nil {
			return err
		}
		tr.Op = function
	}

	tr.Branches = serial.Branches

	return nil
}
