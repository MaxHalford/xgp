package koza

import (
	"encoding/json"

	"github.com/MaxHalford/koza/metrics"
)

// A Task contains information a Program needs to know in order to at least
// make predictions.
type Task struct {
	Metric   metrics.Metric
	NClasses int // Should be equal to 0 if Classification is false
}

func (t Task) binaryClassification() bool {
	return t.Metric.Classification() && t.NClasses == 2
}

func (t Task) multiClassification() bool {
	return t.Metric.Classification() && t.NClasses > 2
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
