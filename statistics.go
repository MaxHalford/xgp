package xgp

import (
	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp/tree"
)

type Statistics struct {
	AvgHeight float64
	AvgNNodes float64
}

func CollectStatistics(GA *gago.GA) Statistics {
	var (
		stats Statistics
		n     float64
	)
	for _, pop := range GA.Populations {
		for _, indi := range pop.Individuals {
			stats.AvgHeight += float64(tree.GetHeight(indi.Genome.(*Program).Root))
			stats.AvgNNodes += float64(tree.GetNNodes(indi.Genome.(*Program).Root))
			n++
		}
	}
	stats.AvgHeight /= n
	stats.AvgNNodes /= n
	return stats
}
