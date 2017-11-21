package koza

import (
	"math"
	"sort"
)

// closestFloat64 performs a binary search to find the index of closest float in
// fs to t.
func closestFloat64(fs []float64, f float64) float64 {
	var i = len(fs) / 2
	if fs[i] > f {
		switch len(fs) {
		case 1:
			return fs[0]
		case 2:
			if f-fs[0] < fs[1]-f {
				return fs[0]
			}
			return fs[1]
		default:
			return closestFloat64(fs[:i+1], f)
		}
	}
	switch len(fs) {
	case 1:
		return fs[0]
	case 2:
		if f-fs[0] < fs[1]-f {
			return fs[0]
		}
		return fs[1]
	default:
		return closestFloat64(fs[i:len(fs)], f)
	}
}

// DynamicRangeSelection implements Dynamic Range Selection (DRS). The
// pseudo-code can be found on page 66 of
// http://goanna.cs.rmit.edu.au/~vc/papers/loveard-phd.pdf.
type DynamicRangeSelection struct {
	cutPoints []float64
	rangeMap  map[float64]float64
}

// Fit DynamicRangeSelection.
func (drs *DynamicRangeSelection) Fit(yTrue, yPred []float64) error {
	// Round each output to it's closest integer
	var mapping = make(map[float64]map[float64]float64)
	for i := range yPred {
		var r = math.Floor(yPred[i] + 0.5)
		if _, ok := mapping[r]; ok {
			mapping[r][yTrue[i]]++
		} else {
			mapping[r] = map[float64]float64{yTrue[i]: 1}
		}
	}
	// Count each class occurence
	drs.rangeMap = make(map[float64]float64)
	for output, classCounts := range mapping {
		var (
			maxCount = math.Inf(-1)
			class    float64
		)
		for c, count := range classCounts {
			if count > maxCount {
				maxCount = count
				class = c
			}
		}
		drs.rangeMap[output] = class
	}
	// Determine cut points
	drs.cutPoints = make([]float64, len(drs.rangeMap))
	var i int
	for cp := range drs.rangeMap {
		drs.cutPoints[i] = cp
		i++
	}
	sort.Float64s(drs.cutPoints)
	return nil
}

// PredictRow DynamicRangeSelection.
func (drs DynamicRangeSelection) PredictRow(y float64) float64 {
	var cutPoint = closestFloat64(drs.cutPoints, y)
	return drs.rangeMap[cutPoint]
}

// Predict DynamicRangeSelection.
func (drs DynamicRangeSelection) Predict(yPred []float64) []float64 {
	var classes = make([]float64, len(yPred))
	for i, y := range yPred {
		classes[i] = drs.PredictRow(y)
	}
	return classes
}

// Clone DynamicRangeSelection.
func (drs DynamicRangeSelection) clone() *DynamicRangeSelection {
	// Copy cut-points
	var cutPoints = make([]float64, len(drs.cutPoints))
	copy(cutPoints, drs.cutPoints)
	// Copy range map
	var rangeMap = make(map[float64]float64)
	for k, v := range drs.rangeMap {
		rangeMap[k] = v
	}
	return &DynamicRangeSelection{
		cutPoints: cutPoints,
		rangeMap:  rangeMap,
	}
}
