package xgp

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDRSJSONEncodeDecode(t *testing.T) {
	var initialDRS = &DynamicRangeSelection{
		cutPoints: []float64{0, 1, 2},
		rangeMap:  map[float64]float64{0: 1, 1: 1, 2: 1},
	}

	// Serialize the initial DRS
	var bytes, err = json.Marshal(initialDRS)
	if err != nil {
		t.Errorf("Expected nil, got '%s'", err.Error())
	}

	// Parse the bytes into a new tree
	var newDRS *DynamicRangeSelection
	err = json.Unmarshal(bytes, &newDRS)
	if err != nil {
		t.Errorf("Expected nil, got '%s'", err.Error())
	}

	// Compare the DRSs
	if !reflect.DeepEqual(newDRS, initialDRS) {
		t.Error("Initial and new DRS do not match")
	}
}
