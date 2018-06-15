package op

import (
	"math/rand"
	"sort"

	"github.com/gonum/floats"
)

// An Operator is a mathematical operator. It has operands that are themselves
// Operators.
type Operator interface {
	Eval(X [][]float64) []float64
	Arity() uint
	Operand(i uint) Operator
	SetOperand(i uint, op Operator) Operator
	Simplify() Operator
	Diff(i uint) Operator
	Name() string
	String() string
}

// walk performs tree traversal on an Operator and applies a given function to
// each suboperator. The given function specifies if the traversal should stop
// or not. Pre-order traversal is used.
func walk(op Operator, f func(op Operator, depth, pos uint) (stop bool)) {
	var (
		step    func(op Operator, depth uint) bool
		counter uint
	)
	step = func(op Operator, depth uint) bool {
		counter++
		// Apply the given function to the Operator
		if f(op, depth, counter-1) {
			return true
		}
		// walk over the operands
		for i := uint(0); i < op.Arity(); i++ {
			if step(op.Operand(i), depth+1) {
				return true
			}
		}
		return false
	}
	step(op, 0)
}

// CalcHeight returns the height of an Operator. The height of an Operator is
// the height of it's root node. The height of a node is the number of edges
// on the longest path between the node and a leaf.
func CalcHeight(op Operator) (d uint) {
	var f = func(op Operator, depth, pos uint) (stop bool) {
		if depth > d {
			d = depth
		}
		return
	}
	walk(op, f)
	return
}

// Count returns the number of Operators that match a condition.
func Count(op Operator, condition func(Operator) bool) (n uint) {
	var f = func(op Operator, depth, pos uint) (stop bool) {
		if condition(op) {
			n++
		}
		return
	}
	walk(op, f)
	return n
}

// CountOps returns the number of Operators in an Operator.
func CountOps(op Operator) uint {
	return Count(op, func(Operator) bool { return true })
}

// CountLeaves returns the number of leaves in an Operator.
func CountLeaves(op Operator) uint {
	return Count(op, func(op Operator) bool { return op.Arity() == 0 })
}

// CountConsts returns the number of Consts in an Operator.
func CountConsts(op Operator) uint {
	return Count(op, func(op Operator) bool { _, ok := op.(Const); return ok })
}

// CountVars returns the number of Vars in an Operator.
func CountVars(op Operator) uint {
	return Count(op, func(op Operator) bool { _, ok := op.(Var); return ok })
}

// Select returns the nth Operator in an Operator.
func Select(op Operator, pos uint) Operator {
	var (
		selected Operator
		f        = func(op Operator, depth, i uint) (stop bool) {
			if i == pos {
				selected = op
				return true
			}
			return false
		}
	)
	walk(op, f)
	return selected
}

// Sample returns a random suboperator.
func Sample(
	op Operator,
	weight func(op Operator, depth uint, rng *rand.Rand) float64,
	rng *rand.Rand,
) (Operator, uint) {
	var (
		weights = make([]float64, CountOps(op))
		f       = func(op Operator, depth, pos uint) (stop bool) {
			weights[pos] = weight(op, depth, rng)
			return
		}
	)
	// Assign weights to each suboperator
	walk(op, f)
	// Calculate the cumulative sum of the weights
	var cs = make([]float64, len(weights))
	floats.CumSum(cs, weights)
	// Sample a random number in [0, cs[-1])
	var r = rng.Float64() * cs[len(cs)-1]
	// Find i where cs[i-1] < r < cs[i]
	var i = uint(sort.SearchFloat64s(cs, r))
	// Return the ith suboperator
	return Select(op, i), i
}

// Replace replaces an Operator if a given condition is met.
func Replace(
	op Operator,
	when func(Operator) bool,
	how func(Operator) Operator,
	stopAtFirst bool,
) Operator {
	var walk func(op Operator) (Operator, bool)
	walk = func(op Operator) (Operator, bool) {
		if when(op) {
			return how(op), stopAtFirst
		}
		for i := uint(0); i < op.Arity(); i++ {
			operand, stop := walk(op.Operand(i))
			op = op.SetOperand(i, operand)
			if stop {
				return op.SetOperand(i, operand), true
			}
		}
		return op, false
	}
	op, _ = walk(op)
	return op
}

// ReplaceAt replaces op's suboperator at position pos.
func ReplaceAt(op Operator, pos uint, with Operator) Operator {
	var when = func(op Operator) bool {
		if pos == 0 {
			return true
		}
		pos--
		return false
	}
	return Replace(op, when, func(op Operator) Operator { return with }, true)
}

// GetConsts returns the Consts contained in an Operator. Pre-order traversal
// is used.
func GetConsts(op Operator) []float64 {
	var (
		values = make([]float64, 0)
		step   = func(op Operator, depth, pos uint) (stop bool) {
			if c, ok := op.(Const); ok {
				values = append(values, c.Value)
			}
			return false
		}
	)
	walk(op, step)
	return values
}

// SetConsts sets the value of each Const. Pre-order traversal is used.
func SetConsts(op Operator, values []float64) Operator {
	var (
		counter uint
		walk    func(op Operator) (Operator, bool)
	)
	walk = func(op Operator) (Operator, bool) {
		if _, ok := op.(Const); ok {
			defer func() { counter++ }()
			return Const{values[counter]}, int(counter) == len(values)
		}
		for i := uint(0); i < op.Arity(); i++ {
			operand, stop := walk(op.Operand(i))
			op = op.SetOperand(i, operand)
			if stop {
				return op, true
			}
		}
		return op, int(counter) == len(values)
	}
	op, _ = walk(op)
	return op
}
