package metrics

import "fmt"

type errMismatchedLengths struct {
	yTrueLen int
	yPredLen int
}

func (e *errMismatchedLengths) Error() string {
	return fmt.Sprintf("Mismatched lengths: len(yTrue) != len(yPred) (%d != %d)", e.yTrueLen, e.yPredLen)
}

type errClassNotFound struct {
	class float64
}

func (e *errClassNotFound) Error() string {
	return fmt.Sprintf("Class not found: %0.f", e.class)
}

type errUnknownMetric struct {
	metricName string
}

func (e *errUnknownMetric) Error() string {
	return fmt.Sprintf("Unknown metric name: '%s'", e.metricName)
}
