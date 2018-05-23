package xgp

import "fmt"

func ExampleConfig() {
	var (
		conf     = NewConfigWithDefaults()
		est, err = conf.NewEstimator()
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(est)
	// Output: Loss metric: mse
	// Evaluation metric: mse
	// Parsimony coefficient: 0
	// Functions: sum,sub,mul,div
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
	// Number of tuning generations: 0
	// Hoist mutation probability: 0.1
	// Subtree mutation probability: 0.1
	// Point mutation probability: 0.1
	// Point mutation rate: 0.3
	// Subtree crossover probability: 0.5
}
