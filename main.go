package main

import (
	"fmt"
	"log"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/tiago/dataframe"
	"github.com/MaxHalford/tiago/metric"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	// Open the training set
	train, err := dataframe.ReadCSV("examples/iris/train.csv", "target", false)
	check(err)
	// Open the test set
	test, err := dataframe.ReadCSV("examples/iris/test.csv", "target", false)
	check(err)

	// Instantiate an Estimator
	var estimator = Estimator{
		DataFrame:       train,
		Metric:          metric.NegativeBinaryF1Score{1},
		Activation:      Sigmoid,
		PVariable:       0.5,
		NodeInitializer: FullNodeInitializer{Height: 3},
		FunctionSet: map[int][]Operator{
			2: []Operator{
				Sum{},
				Difference{},
				Product{},
				Division{},
			},
		},
	}

	// Set the Estimator's GA
	estimator.GA = gago.GA{
		GenomeFactory: estimator.NewProgram,
		NPops:         1,
		PopSize:       100,
		Model: gago.ModGenerational{
			Selector: gago.SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
	}

	// Initialize the Estimator
	estimator.Initialize()

	var generations = 10
	for i := 0; i < generations; i++ {
		estimator.GA.Enhance()
		fmt.Printf("Score: %.3f | %d / %d \n", estimator.GA.Best.Fitness, i+1, generations)
		//for _, pop := range estimator.GA.Populations {
		//	for _, indi := range pop.Individuals {
		//fmt.Println(indi.Genome.(*Program))
		//fmt.Println(indi.Genome.(*Program).PredictDataFrame(train, false))
		//fmt.Println(indi.Genome.(*Program).PredictDataFrame(train, true))
		//fmt.Println(strings.Repeat("-", 30))
		//fmt.Printf("%p\n", indi.Genome.(*Program).Root)
		//}
		//}
		//fmt.Println(strings.Repeat("=", 30))
	}

	yPred := estimator.GA.Best.Genome.(*Program).PredictDataFrame(test, false)
	score, _ := estimator.Metric.Apply(test.Y, yPred)
	fmt.Println(score)
}
