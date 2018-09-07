package meta

import (
	"math"
)

// A LineSearcher finds a good enough step size to a gradient descent problem.
// f is the function we want to minimize given a step size.
type LineSearcher interface {
	Solve(f func(x float64) float64) float64
}

// GoldenLineSearch implements Golden-section line-search.
type GoldenLineSearch struct {
	Min, Max float64 // Initial interval
	Tol      float64
}

// Solve of GoldenLineSearch taken from
// https://www.wikiwand.com/en/Golden-section_search.
func (gls GoldenLineSearch) Solve(f func(x float64) float64) float64 {
	var (
		a, b = gls.Min, gls.Max
		h    = b - a
	)
	if h <= gls.Tol {
		return (a + b) / 2
	}
	var (
		invphi  = (math.Sqrt(5) - 1) / 2
		invphi2 = (3 - math.Sqrt(5)) / 2
		n       = uint(math.Ceil(math.Log(gls.Tol/h) / math.Log(invphi)))
		c       = a + invphi2*h
		d       = a + invphi*h
		yc      = f(c)
		yd      = f(d)
	)
	for i := uint(0); i < n; i++ {
		if yc < yd {
			b = d
			d = c
			yd = yc
			h = invphi * h
			c = a + invphi2*h
			yc = f(c)
		} else {
			a = c
			c = d
			yc = yd
			h = invphi * h
			d = a + invphi*h
			yd = f(d)
		}
	}
	if yc < yd {
		return (a + d) / 2
	}
	return (c + b) / 2
}
