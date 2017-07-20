package main

import (
	"fmt"
	"log"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/tiago/dataframe"
	"github.com/MaxHalford/tiago/gp"
	"github.com/MaxHalford/tiago/metric"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	// Open the training dataframe
	var df, err = dataframe.ReadCSV("examples/gplearn_regression/train.csv", "y", false)
	check(err)

	// Instantiate an Estimator
	var estimator = gp.Estimator{
		DataFrame:       df,
		Metric:          metric.MinkowskiDistance{2},
		Activation:      gp.Identity,
		PVariable:       0.5,
		NodeInitializer: gp.FullNodeInitializer{Height: 3},
		FunctionSet: map[int][]gp.Operator{
			2: []gp.Operator{
				gp.Sum{},
				gp.Difference{},
				gp.Product{},
				gp.Division{},
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

	var generations = 20
	for i := 0; i < generations; i++ {
		estimator.GA.Enhance()
		fmt.Printf("Score: %.3f | %d / %d \n", estimator.GA.Best.Fitness, i+1, generations)
		//for _, pop := range estimator.GA.Populations {
		//	for _, indi := range pop.Individuals {
		//fmt.Println(indi.Genome.(*gp.Program))
		//fmt.Println(indi.Genome.(*gp.Program).PredictDataFrame(df, false))
		//fmt.Println(indi.Genome.(*gp.Program).PredictDataFrame(df, true))
		//fmt.Println(strings.Repeat("-", 30))
		//fmt.Printf("%p\n", indi.Genome.(*gp.Program).Root)
		//}
		//}
		//fmt.Println(strings.Repeat("=", 30))
	}

	fmt.Println(estimator.GA.Best.Genome.(*gp.Program))
	fmt.Println(estimator.GA.Best.Genome.(*gp.Program).PredictDataFrame(df, false))
	fmt.Println(df.Y)
}
