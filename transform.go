package xgp

import "math"

func Identity(x float64) float64 {
	return x
}

func Binary(x float64) float64 {
	if x < 0 {
		return 0
	}
	return 1
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}
