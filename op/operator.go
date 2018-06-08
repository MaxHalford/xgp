package op

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"

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

// CalcHeight returns the height of an Operator. The height of a Operator is
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

// Count returns the number of Operators that match a criteria.
func Count(op Operator, crit func(Operator) bool) (n uint) {
	var f = func(op Operator, depth, pos uint) (stop bool) {
		if crit(op) {
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

// Replace replaces op's suboperator at position pos.
func Replace(op Operator, pos uint, with Operator) Operator {
	var walk func(op Operator) (Operator, bool)
	walk = func(op Operator) (Operator, bool) {
		if pos == 0 {
			return with, true
		}
		pos--
		for i := uint(0); i < op.Arity(); i++ {
			operand, stop := walk(op.Operand(i))
			if stop {
				return op.SetOperand(i, operand), true
			}
		}
		return op, false
	}
	op, _ = walk(op)
	return op
}

type serialOperator struct {
	Type     string           `json:"type"`
	Value    string           `json:"value"`
	Operands []serialOperator `json:"operands"`
}

func serializeOp(op Operator) serialOperator {
	var (
		arity  = op.Arity()
		serial = serialOperator{Operands: make([]serialOperator, arity)}
	)
	for i := uint(0); i < op.Arity(); i++ {
		serial.Operands[i] = serializeOp(op.Operand(i))
	}
	switch op := op.(type) {
	case Const:
		serial.Type = "const"
		serial.Value = strconv.FormatFloat(op.Value, 'f', -1, 64)
	case Var:
		serial.Type = "var"
		serial.Value = strconv.Itoa(int(op.Index))
	default:
		serial.Type = "func"
		serial.Value = op.Name()
	}
	return serial
}

func parseOp(serial serialOperator) (Operator, error) {
	var op Operator
	switch serial.Type {
	case "const":
		val, err := strconv.ParseFloat(serial.Value, 64)
		if err != nil {
			return nil, err
		}
		op = Const{val}
	case "var":
		idx, err := strconv.Atoi(serial.Value)
		if err != nil {
			return nil, err
		}
		op = Var{uint(idx)}
	default:
		function, err := ParseFunc(serial.Value)
		if err != nil {
			return nil, err
		}
		op = function
	}
	// Set the operands; this is where the recursion happens
	for i, serialOperand := range serial.Operands {
		operand, err := parseOp(serialOperand)
		if err != nil {
			return nil, err
		}
		op = op.SetOperand(uint(i), operand)
	}
	return op, nil
}

// MarshalJSON serializes an Operator into JSON.
func MarshalJSON(op Operator) ([]byte, error) {
	return json.Marshal(serializeOp(op))
}

// UnmarshalJSON parses JSON into an Operator.
func UnmarshalJSON(raw []byte) (Operator, error) {
	var serial = serialOperator{}
	if err := json.Unmarshal(raw, &serial); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return parseOp(serial)
}
