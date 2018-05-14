package xgp

import (
	"encoding/json"

	"github.com/MaxHalford/xgp/metrics"
)

// A Task contains information a Program needs to know in order to at least
// make predictions.
type Task struct {
	LossMetric metrics.Metric
	NClasses   int // Should be equal to 0 if the metric is not a classification metric
}

func (t Task) binaryClassification() bool {
	return t.LossMetric.Classification() && t.NClasses == 2
}

func (t Task) multiClassification() bool {
	return t.LossMetric.Classification() && t.NClasses > 2
}

// A serialTask can be serialized and holds information that can be used to
// initialize a Task.
type serialTask struct {
	LossMetricName string `json:"metric_name"`
	NClasses       int    `json:"n_classes"`
}

// serializeTask transforms a Task into a serialTask.
func serializeTask(t Task) (serialTask, error) {
	return serialTask{
		LossMetricName: t.LossMetric.String(),
		NClasses:       t.NClasses,
	}, nil
}

// parseSerialTask recursively transforms a serialTask into a *DynamicRangeSelection.
func parseSerialTask(serial serialTask) (*Task, error) {
	var metric, err = metrics.GetMetric(serial.LossMetricName, 1)
	if err != nil {
		return nil, err
	}
	return &Task{
		LossMetric: metric,
		NClasses:   serial.NClasses,
	}, nil
}

// MarshalJSON serializes a Task into JSON bytes. A serialTask is used as an
// intermediary.
func (t Task) MarshalJSON() ([]byte, error) {
	var serial, err = serializeTask(t)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a *Task. A serialTask is used as an
// intermediary.
func (t *Task) UnmarshalJSON(bytes []byte) error {
	var serial serialTask
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedTask, err = parseSerialTask(serial)
	if err != nil {
		return err
	}
	*t = *parsedTask
	return nil
}
