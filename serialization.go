package koza

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/MaxHalford/koza/metrics"
)

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

// A serialTask can be serialized and holds information that can be used to
// initialize a Task.
type serialTask struct {
	MetricName string `json:"metric_name"`
	NClasses   int    `json:"n_classes"`
}

// serializeTask transforms a Task into a serialTask.
func serializeTask(task Task) (serialTask, error) {
	return serialTask{
		MetricName: task.Metric.String(),
		NClasses:   task.NClasses,
	}, nil
}

// parseSerialTask recursively transforms a serialTask into a *DynamicRangeSelection.
func parseSerialTask(serial serialTask) (*Task, error) {
	var metric, err = metrics.GetMetric(serial.MetricName, 1)
	if err != nil {
		return nil, err
	}
	return &Task{
		Metric:   metric,
		NClasses: serial.NClasses,
	}, nil
}

// MarshalJSON serializes a Task into JSON bytes. A serialTask is used as an
// intermediary.
func (task Task) MarshalJSON() ([]byte, error) {
	var serial, err = serializeTask(task)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a *Task. A serialTask is used as an
// intermediary.
func (task *Task) UnmarshalJSON(bytes []byte) error {
	var serial serialTask
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedTask, err = parseSerialTask(serial)
	if err != nil {
		return err
	}
	*task = *parsedTask
	return nil
}

// SaveProgramToJSON saves a Program to a JSON file.
func SaveProgramToJSON(program Program, path string) error {
	var bytes, err = json.Marshal(&program)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadProgramFromJSON loads a Program from a JSON file.
func LoadProgramFromJSON(path string) (Program, error) {
	var (
		program    Program
		bytes, err = ioutil.ReadFile(path)
	)
	if err != nil {
		return program, err
	}
	err = json.Unmarshal(bytes, &program)
	if err != nil {
		return program, err
	}
	return program, nil
}
