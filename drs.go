package xgp

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
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

// PredictPartial DynamicRangeSelection.
func (drs DynamicRangeSelection) PredictPartial(y float64) float64 {
	var cutPoint = closestFloat64(drs.cutPoints, y)
	return drs.rangeMap[cutPoint]
}

// Predict DynamicRangeSelection.
func (drs DynamicRangeSelection) Predict(yPred []float64) []float64 {
	var classes = make([]float64, len(yPred))
	for i, y := range yPred {
		classes[i] = drs.PredictPartial(y)
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

// A serialDRS can be serialized and holds information that can be used to
// initialize a DynamicRangeSelection.
type serialDRS struct {
	CutPoints []float64         `json:"cut_points"`
	RangeMap  map[string]string `json:"range_map"`
}

// serializeDRS transforms a *DynamicRangeSelection into a serialDRS.
func serializeDRS(drs *DynamicRangeSelection) (serialDRS, error) {
	var serial = serialDRS{
		CutPoints: drs.cutPoints,
		RangeMap:  make(map[string]string),
	}
	for k, v := range drs.rangeMap {
		var (
			ks = strconv.FormatFloat(k, 'f', -1, 64)
			vs = strconv.FormatFloat(v, 'f', -1, 64)
		)
		serial.RangeMap[ks] = vs
	}
	return serial, nil
}

// parseSerialDRS recursively transforms a serialDRS into a *DynamicRangeSelection.
func parseSerialDRS(serial serialDRS) (*DynamicRangeSelection, error) {
	var drs = &DynamicRangeSelection{
		cutPoints: serial.CutPoints,
		rangeMap:  make(map[float64]float64),
	}
	for k, v := range serial.RangeMap {
		kf, err := strconv.ParseFloat(k, 64)
		if err != nil {
			return nil, err
		}
		vf, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		drs.rangeMap[kf] = vf
	}
	return drs, nil
}

// MarshalJSON serializes a *DynamicRangeSelection into JSON bytes. A serialDRS
// is used as an intermediary.
func (drs *DynamicRangeSelection) MarshalJSON() ([]byte, error) {
	var serial, err = serializeDRS(drs)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a *DynamicRangeSelection. A serialDRS is
// used as an intermediary.
func (drs *DynamicRangeSelection) UnmarshalJSON(bytes []byte) error {
	var serial serialDRS
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedDRS, err = parseSerialDRS(serial)
	if err != nil {
		return err
	}
	*drs = *parsedDRS
	return nil
}
