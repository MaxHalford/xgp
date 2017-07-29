package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataframe"
	"github.com/MaxHalford/xgp/metric"
	"github.com/MaxHalford/xgp/tree"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "xgp"
	app.Usage = "Genetic programming for machine learning tasks"
	app.Compiled = time.Now()

	app.Commands = []cli.Command{
		{
			Name:  "fit",
			Usage: "Fits a model on a dataset",
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
					Value: "mse",
					Usage: "What kind of metric to use",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "model.json",
					Usage: "Path for the output model",
				},
			},
			Action: func(c *cli.Context) error {
				var file = c.Args().First()
				// Check if the file exists
				if err := fileExists(file); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				// Determine the task to perform
				isClassification := c.String("task") == "classification"
				// Load the training set in memory
				train, err := dataframe.ReadCSV(file, c.String("target_col"), isClassification)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				// Instantiate an Estimator
				estimator := xgp.Estimator{
					DataFrame:       train,
					Metric:          metric.MeanSquaredError{},
					Transform:       xgp.Identity,
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

				// Save the best Program's root Node
				var bestProg = estimator.GA.Best.Genome.(*xgp.Program)
				xgp.SaveNodeToJSON(bestProg.Root, c.String("output"))

				return nil

			},
		},
		{
			Name:  "predict",
			Usage: "Predicts a dataset with a trained model",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "model, m",
					Value: "model.json",
					Usage: "Filename of the trained model",
				},
				cli.StringFlag{
					Name:  "target_col, tc",
					Value: "target",
					Usage: "Name of the target column",
				},
			},
			Action: func(c *cli.Context) error {
				var file = c.Args().First()
				// Check if the file exists
				if err := fileExists(file); err != nil {
					return err
				}
				// Load the test set in memory
				test, err := dataframe.ReadCSV(file, c.String("target_col"), false)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				// Load the model
				node, err := xgp.LoadNodeFromJSON(c.String("model"))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				var prog = xgp.Program{
					Root: node,
					Estimator: &xgp.Estimator{
						Metric:    metric.MeanSquaredError{},
						Transform: xgp.Identity,
					},
				}

				var yPred = prog.PredictDataFrame(test, prog.Estimator.Transform)
				var score, _ = prog.Estimator.Metric.Apply(test.Y, yPred)
				fmt.Printf("Test score: %.3f\n", score)

				return nil
			},
		},
		{
			Name:  "todot",
			Usage: "Creates a .dot file from a model",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "model, m",
					Value: "model.json",
					Usage: "Filename of the trained model",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "model.dot",
					Usage: "Path for the output file",
				},
			},
			Action: func(c *cli.Context) error {
				// Load the model
				node, err := xgp.LoadNodeFromJSON(c.String("model"))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				// Create the output file
				file, err := os.Create(c.String("output"))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				defer file.Close()
				// Make the Graphviz representation
				var graphviz = tree.GraphvizDisplay{}
				file.WriteString(graphviz.Apply(node))
				file.Sync()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
