package main

import (
	"fmt"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataframe"
	"github.com/urfave/cli"
)

var fitCmd = cli.Command{
	Name:  "fit",
	Usage: "Fits a program on a dataset",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "target_col, tc",
			Value: "target",
			Usage: "Name of the target column",
		},
		cli.IntFlag{
			Name:  "generations, g",
			Value: 10,
			Usage: "Number of generations",
		},
		cli.StringFlag{
			Name:  "task, t",
			Value: "regression",
			Usage: "What kind of task to perform",
		},
		cli.StringFlag{
			Name:  "metric, m",
			Value: "mean_squared_error",
			Usage: "What kind of metric to use",
		},
		cli.IntFlag{
			Name:  "class, c",
			Value: 1,
			Usage: "Which class to apply the metric to if applicable",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "program.json",
			Usage: "Path for the output program",
		},
	},
	Action: func(c *cli.Context) error {
		// Check if the training file exists
		var file = c.Args().First()
		if err := fileExists(file); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		// Determine the task to perform
		isClassification := c.String("task") == "classification"

		// Determine the metric to use
		metric, err := getMetric(c.String("metric"), c.Float64("class"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		// Load the training set in memory
		train, err := dataframe.ReadCSV(file, c.String("target_col"), isClassification)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		// Instantiate an Estimator
		estimator := xgp.Estimator{
			DataFrame:       train,
			Metric:          metric,
			Transform:       xgp.Identity{},
			PVariable:       0.5,
			NodeInitializer: xgp.FullNodeInitializer{Height: 3},
			FunctionSet: map[int][]xgp.Operator{
				2: []xgp.Operator{
					xgp.Sum{},
					xgp.Difference{},
					xgp.Product{},
					xgp.Division{},
				},
			},
		}
		// Instantiate a GA
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

		for i := 0; i < c.Int("generations"); i++ {
			estimator.GA.Enhance()
			fmt.Printf("Score: %.3f | %d / %d \n", estimator.GA.Best.Fitness, i+1, c.Int("generations"))
		}

		// Save the best Program
		var bestProg = estimator.GA.Best.Genome.(*xgp.Program)
		xgp.SaveProgramToJSON(*bestProg, c.String("output"))

		return nil

	},
}
