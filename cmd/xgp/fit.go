package main

import (
	"math"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataframe"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/urfave/cli"
)

func monitorProgress(ch <-chan float64, done chan bool) {
	uiprogress.Start()
	var (
		bar         = uiprogress.AddBar(cap(ch))
		green       = color.New(color.FgGreen).SprintfFunc()
		bestFitness = math.Inf(1)
	)
	bar.AppendCompleted()
	bar.PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return green("Best fitness => %.5f", bestFitness)
	})
	for i := 0; i < cap(ch); i++ {
		var fitness = <-ch
		bar.Incr()
		if fitness < bestFitness {
			bestFitness = fitness
		}
	}
	uiprogress.Stop()
	done <- true
}

func exitCLI(err error) *cli.ExitError {
	var red = color.New(color.FgRed).SprintFunc()
	return cli.NewExitError(red(err.Error()), 1)
}

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
		cli.IntFlag{
			Name:  "tuningGenerations, tg",
			Value: 5,
			Usage: "Number of tuning generations",
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
			return exitCLI(err)
		}

		// Determine the task to perform
		isClassification := c.String("task") == "classification"

		// Determine the metric to use
		metric, err := getMetric(c.String("metric"), c.Float64("class"))
		if err != nil {
			return exitCLI(err)
		}

		// Instantiate an Estimator
		estimator := xgp.Estimator{
			Metric:    metric,
			Transform: xgp.Identity{},
			PVariable: 0.5,
			NodeInitializer: xgp.RampedHaldAndHalfInitializer{
				MinHeight: 2,
				MaxHeight: 5,
				PLeaf:     0.5,
			},
			FunctionSet: map[int][]xgp.Operator{
				2: []xgp.Operator{
					xgp.Sum{},
					xgp.Difference{},
					xgp.Product{},
					xgp.Division{},
				},
			},
		}

		// Set the initial GA
		estimator.GA = &gago.GA{
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

		// Set the tuning GA
		estimator.TuningGA = &gago.GA{
			GenomeFactory: estimator.NewProgramTuner,
			NPops:         1,
			PopSize:       20,
			Model: gago.ModGenerational{
				Selector: gago.SelTournament{
					NContestants: 3,
				},
				MutRate: 0.5,
			},
		}

		// Load the training set in memory
		df, err := dataframe.ReadCSV(file, c.String("target_col"), isClassification)
		if err != nil {
			return exitCLI(err)
		}

		// Monitor progress
		color.Blue(
			"Fitting for %d generations and tuning for %d generations",
			c.Int("generations"),
			c.Int("tuningGenerations"),
		)
		var (
			ch   = make(chan float64, c.Int("generations")+c.Int("tuningGenerations"))
			done = make(chan bool)
		)
		go monitorProgress(ch, done)

		// Fit the Estimator
		err = estimator.Fit(df, c.Int("generations"), ch)
		if err != nil {
			return exitCLI(err)
		}

		// Tune the Estimator
		err = estimator.Tune(df, c.Int("tuningGenerations"), ch)
		if err != nil {
			return exitCLI(err)
		}

		// Save the best Program
		<-done
		color.Blue("Saving best program to '%s'...", c.String("output"))
		var bestProg, _ = estimator.BestProgram()
		err = xgp.SaveProgramToJSON(*bestProg, c.String("output"))
		if err != nil {
			return exitCLI(err)
		}

		return nil

	},
}
