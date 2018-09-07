package meta

import (
	"math"
	"testing"
)

func TestGoldenLineSearch(t *testing.T) {
	var (
		f = func(x float64) float64 {
			return math.Pow((x - 2), x)
		}
		gls = GoldenLineSearch{
			Min: 1,
			Max: 5,
			Tol: 1e-5,
		}
		step = gls.Solve(f)
	)
	if math.Abs(step-2) > 1e-5 {
		t.Errorf("Expected 2, got %f", step)
	}
}

/*func TestBacktracking(t *testing.T) {
	var (
		f = func(x []float64) float64 {
			var y float64
			for _, xi := range x {
				y += xi * xi
			}
			return y
		}
		x  = []float64{5, 5}
		g  = []float64{2 * 5, 2 * 5}
		bt = Backtracking{
			InitStep:    1,
			MinStep:     0,
			Decrease:    0.5,
			Contraction: 0.5,
			MaxIter:     100,
		}
	)
	var step, xx = bt.Solve(f, x, g)
	if step != 0.5 {
		t.Errorf("Expected 0, got %f", step)
	}
	if !reflect.DeepEqual(xx, []float64{0, 0}) {
		t.Errorf("Expected %v, got %v", []float64{0, 0}, xx)
	}
}
*/
