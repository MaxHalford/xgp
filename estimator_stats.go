package xgp

// AvgProgramHeight returns the average Program height.
func (est Estimator) AvgProgramHeight() float64 {
	var (
		total int
		n     float64
	)
	for _, pop := range est.GA.Populations {
		for _, indi := range pop.Individuals {
			total += indi.Genome.(*Program).Tree.Height()
			n++
		}
	}
	return float64(total) / n
}
