package tree

// A Cache stores results produced by Trees.
type Cache struct {
	results map[string][]float64
}

func (c Cache) Get(key string) []float64 {
	if c.results == nil {
		return nil
	}
	if y, ok := c.results[key]; ok {
		return y
	}
	return nil
}

func (c *Cache) Set(key string, y []float64) {
	c.results[key] = y
}

func (c Cache) Keys() []string {
	var (
		ks = make([]string, len(c.results))
		i  int
	)
	for k := range c.results {
		ks[i] = k
		i++
	}
	return ks
}

// NewCache returns a new Cache.
func NewCache() *Cache {
	var c Cache
	c.results = make(map[string][]float64)
	return &c
}
