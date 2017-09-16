package xgp

import "fmt"

type errUnknownMetric struct {
	metricName string
}

func (e *errUnknownMetric) Error() string {
	return fmt.Sprintf("Unknown metric name: %s", e.metricName)
}
