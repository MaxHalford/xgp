package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataframe"
	"github.com/MaxHalford/xgp/metric"
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
					Value: 30,
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
					Value: "model.xgp",
					Usage: "Name of the output model",
				},
			},
			Action: func(c *cli.Context) error {
				file := c.Args().First()
				fmt.Println(file)
				// Check if the file exists
				if err := fileExists(file); err != nil {
					return err
				}
				// Determine the task to perform
				isClassification := c.String("task") == "classification"
				// Load the training set in memory
				train, err := dataframe.ReadCSV(file, c.String("target_col"), isClassification)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				// Instantiate an Estimator
				log.Println(c.String("metric"))
				estimator := xgp.Estimator{
					DataFrame:       train,
					Metric:          metric.NegativeBinaryF1Score{},
					Transform:       xgp.Sigmoid,
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

				return nil

			},
		},
		{
			Name:  "predict",
			Usage: "Predicts a dataset with a trained model",
			Action: func(c *cli.Context) error {
				// yPred := estimator.GA.Best.Genome.(*xgp.Program).PredictDataFrame(test, false)
				// score, _ := estimator.Metric.Apply(test.Y, yPred)
				// fmt.Println(score)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
