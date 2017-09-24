package xgp

import (
	"math"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataset"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(fitCmd)

	fitCmd.Flags().IntVarP(&class, "class", "c", 1, "Which class to apply the metric to if applicable")
	fitCmd.Flags().StringVarP(&funcsString, "functions", "f", "+,-,*,/", "Allowed functions")
	fitCmd.Flags().IntVarP(&generations, "generations", "g", 30, "Number of generations")
	fitCmd.Flags().IntVarP(&maxHeight, "max_height", "u", 6, "Max program height used in ramped half-and-half initialization")
	fitCmd.Flags().IntVarP(&minHeight, "min_height", "l", 2, "Min program height used in ramped half-and-half initialization")
	fitCmd.Flags().StringVarP(&metricName, "metric", "m", "mse", "Metric to use, this determines if the task is classification or regression")
	fitCmd.Flags().StringVarP(&outputName, "output", "o", "program.json", "Path where to save the output program")
	fitCmd.Flags().Float64VarP(&pLeaf, "p_leaf", "p", 0.5, "Probability of generating a leaf node in ramped half-and-half initialization")
	fitCmd.Flags().Float64VarP(&pVariable, "p_variable", "v", 0.5, "Probability of picking a variable and not a constant when generating leaf nodes")
	fitCmd.Flags().StringVarP(&targetCol, "target_col", "y", "y", "Name of the target column in the training set")
	fitCmd.Flags().IntVarP(&tuningGenerations, "t_generations", "t", 30, "Number of tuning generations")
}

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

var fitCmd = &cobra.Command{
	Use:   "fit",
	Short: "Fits a program to a dataset",
	Long:  "Fits a program to a dataset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if the training file exists
		var file = args[0]
		if err := fileExists(file); err != nil {
			return err
		}

		// Determine the metric to use
		metric, err := metrics.GetMetric(metricName, class)
		if err != nil {
			return err
		}

		// Determine the functions to use
		functions, err := parseStringFuncs(funcsString)
		if err != nil {
			return err
		}

		// Load the training set in memory
		dataset, err := dataset.ReadCSV(file, targetCol, metric.Classification())
		if err != nil {
			return err
		}

		// Instantiate an Estimator
		estimator := xgp.Estimator{
			Metric:    metric,
			Transform: xgp.Identity{},
			PVariable: pVariable,
			NodeInitializer: xgp.RampedHaldAndHalfInitializer{
				MinHeight: minHeight,
				MaxHeight: maxHeight,
				PLeaf:     pLeaf,
			},
			Functions:         functions,
			Generations:       generations,
			TuningGenerations: tuningGenerations,
			ProgressChan:      make(chan float64, generations+tuningGenerations),
		}

		// Set the initial GA
		estimator.GA = &gago.GA{
			NewGenome: estimator.NewProgram,
			NPops:     1,
			PopSize:   100,
			Model: gago.ModGenerational{
				Selector: gago.SelTournament{
					NContestants: 3,
				},
				MutRate: 0.5,
			},
		}

		// Set the tuning GA
		estimator.TuningGA = &gago.GA{
			NewGenome: estimator.NewProgramTuner,
			NPops:     1,
			PopSize:   20,
			Model: gago.ModGenerational{
				Selector: gago.SelTournament{
					NContestants: 3,
				},
				MutRate: 0.5,
			},
		}

		// Monitor progress
		color.Blue(
			"Fitting for %d generations and tuning for %d generations",
			generations,
			tuningGenerations,
		)
		var done = make(chan bool)
		go monitorProgress(estimator.ProgressChan, done)

		// Fit the Estimator
		err = estimator.Fit(dataset.X, dataset.Y)
		if err != nil {
			return err
		}

		// Save the best Program
		<-done
		color.Blue("Saving best program to '%s'...", outputName)
		var bestProg, _ = estimator.BestProgram()
		err = xgp.SaveProgramToJSON(*bestProg, outputName)
		if err != nil {
			return err
		}

		return nil
	},
}
