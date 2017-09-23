package dataset

type ClassMap struct {
	Map        map[string]float64
	ReverseMap map[float64]string
	N          float64
}

func (cm *ClassMap) Get(c string) float64 {
	if v, ok := cm.Map[c]; ok {
		return v
	}
	cm.Map[c] = cm.N
	cm.ReverseMap[cm.N] = c
	cm.N++
	return cm.Map[c]
}

func MakeClassMap() *ClassMap {
	return &ClassMap{
		Map:        make(map[string]float64),
		ReverseMap: make(map[float64]string),
	}
}
