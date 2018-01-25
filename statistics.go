package koza

import (
	"github.com/MaxHalford/gago"
)

type statistics struct {
	avgHeight     float64
	avgNOperators float64
}

func collectStats(GA *gago.GA) statistics {
	var (
		stats statistics
		n     float64
	)
	for _, pop := range GA.Populations {
		for _, indi := range pop.Individuals {
			stats.avgHeight += float64(indi.Genome.(*Program).Tree.Height())
			stats.avgNOperators += float64(indi.Genome.(*Program).Tree.Size())
			n++
		}
	}
	stats.avgHeight /= n
	stats.avgNOperators /= n
	return stats
}
