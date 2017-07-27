package main

import "github.com/MaxHalford/tiago/tree"

func (est Estimator) AvgProgramHeight() float64 {
	var (
		total int
		n     float64
	)
	for _, pop := range est.GA.Populations {
		for _, indi := range pop.Individuals {
			total += tree.GetHeight(indi.Genome.(*Program).Root)
			n++
		}
	}
	return float64(total) / n
}
