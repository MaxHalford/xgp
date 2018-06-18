package op

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestEval(t *testing.T) {
	var testCases = []struct {
		in  [][]float64
		op  Operator
		out []float64
	}{
		{
			in:  [][]float64{[]float64{1, -1, 0}},
			op:  Abs{Var{0}},
			out: []float64{1, 1, 0},
		},
		{
			in:  [][]float64{[]float64{1, -1, 0}},
			op:  Const{42},
			out: []float64{42, 42, 42},
		},
		{
			in:  [][]float64{[]float64{-1, 0, 1}},
			op:  Add{Var{0}, Const{1}},
			out: []float64{0, 1, 2},
		},
		{
			in:  [][]float64{[]float64{0, math.Pi}},
			op:  Cos{Var{0}},
			out: []float64{1, -1},
		},
		{
			in: [][]float64{
				[]float64{0, 1, 2},
				[]float64{1, 0, 1},
			},
			op:  Div{Var{0}, Var{1}},
			out: []float64{0, 1, 2},
		},
		{
			in:  [][]float64{[]float64{0, 1, 2}},
			op:  Inv{Var{0}},
			out: []float64{1, 1, 0.5},
		},
		{
			in:  [][]float64{[]float64{0, 1, 2}},
			op:  Max{Var{0}, Const{1}},
			out: []float64{1, 1, 2},
		},
		{
			in:  [][]float64{[]float64{0, 1, 2}},
			op:  Min{Var{0}, Const{1}},
			out: []float64{0, 1, 1},
		},
		{
			in:  [][]float64{[]float64{0, 1, 2}},
			op:  Mul{Var{0}, Const{-1}},
			out: []float64{0, -1, -2},
		},
		{
			in:  [][]float64{[]float64{0, 1, 2}},
			op:  Neg{Var{0}},
			out: []float64{0, -1, -2},
		},
		{
			in:  [][]float64{[]float64{0, 0.5 * math.Pi}},
			op:  Sin{Var{0}},
			out: []float64{0, 1},
		},
		{
			in:  [][]float64{[]float64{-2, 1, 2}},
			op:  Square{Var{0}},
			out: []float64{4, 1, 4},
		},
		{
			in:  [][]float64{[]float64{0, 1, 2}},
			op:  Sub{Var{0}, Const{-1}},
			out: []float64{1, 2, 3},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := tc.op.Eval(tc.in)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("Expected %v, got %v", tc.out, out)
			}
		})
	}
}

func TestArity(t *testing.T) {
	var testCases = []struct {
		in  Operator
		out uint
	}{
		{Const{42}, 0},
		{Var{42}, 0},
		{Abs{}, 1},
		{Add{}, 2},
		{Cos{}, 1},
		{Div{}, 2},
		{Inv{}, 1},
		{Max{}, 2},
		{Min{}, 2},
		{Sin{}, 1},
		{Square{}, 1},
		{Sub{}, 2},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := tc.in.Arity()
			if out != tc.out {
				t.Errorf("Expected %d, got %d", tc.out, out)
			}
		})
	}
}

func TestString(t *testing.T) {
	var testCases = []struct {
		in  Operator
		out string
	}{
		{
			in:  Const{-42.42},
			out: "-42.42",
		},
		{
			in:  Var{42},
			out: "x42",
		},
		{
			in:  Abs{Var{0}},
			out: "|x0|",
		},
		{
			in:  Add{Var{0}, Add{Var{1}, Const{42}}},
			out: "x0+x1+42",
		},
		{
			in:  Cos{Add{Var{0}, Var{1}}},
			out: "cos(x0+x1)",
		},
		{
			in:  Div{Add{Var{0}, Var{1}}, Const{5}},
			out: "(x0+x1)/5",
		},
		{
			in:  Inv{Add{Var{0}, Var{1}}},
			out: "1/(x0+x1)",
		},
		{
			in:  Max{Add{Var{0}, Var{1}}, Const{3}},
			out: "max(x0+x1, 3)",
		},
		{
			in:  Min{Add{Var{0}, Var{1}}, Const{3}},
			out: "min(x0+x1, 3)",
		},
		{
			in:  Sin{Max{Var{0}, Var{1}}},
			out: "sin(max(x0, x1))",
		},
		{
			in:  Square{Max{Var{0}, Var{1}}},
			out: "(max(x0, x1))²",
		},
		{
			in:  Sub{Max{Var{0}, Var{1}}, Var{2}},
			out: "max(x0, x1)-x2",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := tc.in.String()
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}

func TestHeight(t *testing.T) {
	var testCases = []struct {
		op Operator
		d  uint
	}{
		{
			op: Const{42},
			d:  0,
		},
		{
			op: Cos{Const{42}},
			d:  1,
		},
		{
			op: Mul{Cos{Const{42}}, Cos{Const{42}}},
			d:  2,
		},
		{
			op: Mul{Cos{Const{42}}, Cos{Cos{Const{42}}}},
			d:  3,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			d := CalcHeight(tc.op)
			if d != tc.d {
				t.Errorf("Expected %d, got %d", tc.d, d)
			}
		})
	}
}

func TestCountOps(t *testing.T) {
	var testCases = []struct {
		op Operator
		n  uint
	}{
		{
			op: Const{42},
			n:  1,
		},
		{
			op: Cos{Const{42}},
			n:  2,
		},
		{
			op: Mul{Cos{Const{42}}, Cos{Const{42}}},
			n:  5,
		},
		{
			op: Mul{Cos{Const{42}}, Cos{Cos{Const{42}}}},
			n:  6,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			n := CountOps(tc.op)
			if n != tc.n {
				t.Errorf("Expected %d, got %d", tc.n, n)
			}
		})
	}
}

func TestSelect(t *testing.T) {
	var testCases = []struct {
		in  Operator
		pos uint
		out Operator
	}{
		{
			in:  Const{42},
			pos: 0,
			out: Const{42},
		},
		{
			in:  Cos{Const{42}},
			pos: 0,
			out: Cos{Const{42}},
		},
		{
			in:  Cos{Const{42}},
			pos: 1,
			out: Const{42},
		},
		{
			in:  Add{Const{42}, Var{1}},
			pos: 0,
			out: Add{Const{42}, Var{1}},
		},
		{
			in:  Add{Const{42}, Var{1}},
			pos: 1,
			out: Const{42},
		},
		{
			in:  Add{Const{42}, Var{1}},
			pos: 2,
			out: Var{1},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := Select(tc.in, tc.pos)
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}

func TestSample(t *testing.T) {
	var (
		rng       = rand.New(rand.NewSource(time.Now().UnixNano()))
		testCases = []struct {
			in     Operator
			weight func(op Operator, depth uint, rng *rand.Rand) float64
			out    Operator
		}{
			{
				in: Const{42},
				weight: func(op Operator, depth uint, rng *rand.Rand) float64 {
					return 1
				},
				out: Const{42},
			},
			{
				in: Cos{Const{42}},
				weight: func(op Operator, depth uint, rng *rand.Rand) float64 {
					switch op.(type) {
					case Const:
						return 1
					default:
						return 0
					}
				},
				out: Const{42},
			},
			{
				in: Add{Var{1}, Cos{Const{42}}},
				weight: func(op Operator, depth uint, rng *rand.Rand) float64 {
					switch depth {
					case 2:
						return 1
					default:
						return 0
					}
				},
				out: Const{42},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out, _ := Sample(tc.in, tc.weight, rng)
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}

func TestReplaceAt(t *testing.T) {
	var testCases = []struct {
		in   Operator
		pos  uint
		with Operator
		out  Operator
	}{
		{
			in:   Add{Const{0}, Var{0}},
			pos:  0,
			with: Const{42},
			out:  Const{42},
		},
		{
			in:   Add{Const{0}, Var{0}},
			pos:  1,
			with: Const{42},
			out:  Add{Const{42}, Var{0}},
		},
		{
			in:   Add{Const{0}, Var{0}},
			pos:  2,
			with: Const{42},
			out:  Add{Const{0}, Const{42}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := ReplaceAt(tc.in, tc.pos, tc.with)
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	var (
		op         = Add{Var{0}, Const{42.24}}
		bytes, err = MarshalJSON(op)
	)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
		return
	}
	newOp, err := UnmarshalJSON(bytes)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
		return
	}
	if newOp != op {
		t.Errorf("Expected %s, got %s", op, newOp)
		return
	}
}

func TestGetConsts(t *testing.T) {
	var testCases = []struct {
		op     Operator
		values []float64
	}{
		{
			op:     Const{42},
			values: []float64{42},
		},
		{
			op:     Add{Const{1}, Const{0}},
			values: []float64{1, 0},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			values := GetConsts(tc.op)
			if !reflect.DeepEqual(values, tc.values) {
				t.Errorf("Expected %v, got %v", tc.values, values)
			}
		})
	}
}

func TestSetConsts(t *testing.T) {
	var testCases = []struct {
		in     Operator
		values []float64
		out    Operator
	}{
		{
			in:     Const{42},
			values: []float64{43},
			out:    Const{43},
		},
		{
			in:     Add{Const{1}, Const{0}},
			values: []float64{0, 1},
			out:    Add{Const{0}, Const{1}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := SetConsts(tc.in, tc.values)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("Expected %v, got %v", tc.out, out)
			}
		})
	}
}
