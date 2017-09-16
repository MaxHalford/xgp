package main

import (
	"fmt"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataframe"
	"github.com/urfave/cli"
)

var predictCmd = cli.Command{
	Name:  "predict",
	Usage: "Predicts a dataset with a program",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "program, p",
			Value: "program.json",
			Usage: "Path to the trained program",
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
			Name:  "target_col, tc",
			Value: "target",
			Usage: "Name of the target column",
		},
	},
	Action: func(c *cli.Context) error {
		// Check if the file exists
		var file = c.Args().First()
		if err := fileExists(file); err != nil {
			return err
		}

		// Determine the metric to use
		metric, err := getMetric(c.String("metric"), c.Float64("class"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		// Load the test set in memory
		test, err := dataframe.ReadCSV(file, c.String("target_col"), false)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		// Load the program
		prog, err := xgp.LoadProgramFromJSON(c.String("program"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		var yPred = prog.PredictDataFrame(test)
		var score, _ = metric.Apply(test.Y, yPred, nil)
		fmt.Printf("Test score: %.3f\n", score)

		return nil
	},
}
