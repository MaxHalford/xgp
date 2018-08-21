package meta

import (
	"fmt"
	"math"
	"time"
)

func argmin(fs []float64) uint {
	var (
		i   int
		min = math.Inf(1)
	)
	for j, f := range fs {
		if f < min {
			i = j
			min = f
		}
	}
	return uint(i)
}

// mean computes the mean of a float64 slice.
func mean(x []float64) (m float64) {
	for _, xi := range x {
		m += xi
	}
	return m / float64(len(x))
}

// fmtDuration the "hours:minutes:seconds" representation of a time.Duration.
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
