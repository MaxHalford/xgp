package dataset

type classMap struct {
	Map        map[string]float64
	ReverseMap map[float64]string
	N          int
}

func (cm *classMap) Get(c string) float64 {
	if v, ok := cm.Map[c]; ok {
		return v
	}
	cm.N++
	cm.Map[c] = float64(cm.N)
	cm.ReverseMap[float64(cm.N)] = c
	return cm.Map[c]
}

func makeClassMap() *classMap {
	return &classMap{
		Map:        make(map[string]float64),
		ReverseMap: make(map[float64]string),
	}
}
