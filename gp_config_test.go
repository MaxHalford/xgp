package xgp

import (
	"fmt"

	"github.com/MaxHalford/xgp/metrics"
)

func ExampleGPConfig() {
	var conf = NewDefaultGPConfig()
	conf.LossMetric = metrics.F1{}
	var est, err = conf.NewGP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(est)
	// Output:
	// Loss metric: neg_f1
	// Evaluation metric: f1
	// Parsimony coefficient: 0
	// Polish best program: true
	// Functions: add,sub,mul,div
	// Constant minimum: -5
	// Constant maximum: 5
	// Constant probability: 0.5
	// Full initialization probability: 0.5
	// Terminal probability: 0.3
	// Minimum height: 3
	// Maximum height: 5
	// Number of populations: 1
	// Number of individuals per population: 100
	// Number of generations: 30
	// Hoist mutation probability: 0.1
	// Subtree mutation probability: 0.1
	// Point mutation probability: 0.1
	// Point mutation rate: 0.3
	// Subtree crossover probability: 0.5
}
