package xgp

import cache "github.com/patrickmn/go-cache"

// A NodeCache can store Node results.
type NodeCache struct {
	c *cache.Cache
}

// Set Node results.
func (pc *NodeCache) Set(node *Node, results []float64) {
	pc.c.Set(node.String(), results, cache.DefaultExpiration)
}

// Get Node results.
func (pc NodeCache) Get(node *Node) ([]float64, bool) {
	var results, found = pc.c.Get(node.String())
	if found {
		return results.([]float64), found
	}
	return nil, found
}
