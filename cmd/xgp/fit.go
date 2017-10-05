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

var (
	fitClass             int
	fitFuncsString       string
	fitGenerations       int
	fitMaxHeight         int
	fitMetricName        string
	fitMinHeight         int
	fitOutputName        string
	fitPLeaf             float64
	fitPVariable         float64
	fitRounds            int
	fitTargetCol         string
	fitTuningGenerations int
)

func init() {
	RootCmd.AddCommand(fitCmd)

	fitCmd.Flags().IntVarP(&fitClass, "class", "c", 1, "Which class to apply the metric to if applicable")
	fitCmd.Flags().StringVarP(&fitFuncsString, "functions", "f", "+,-,*,/", "Allowed functions")
	fitCmd.Flags().IntVarP(&fitGenerations, "generations", "g", 30, "Number of generations")
	fitCmd.Flags().IntVarP(&fitMaxHeight, "max_height", "u", 6, "Max program height used in ramped half-and-half initialization")
	fitCmd.Flags().StringVarP(&fitMetricName, "metric", "m", "mae", "Metric to use, this determines if the task is classification or regression")
	fitCmd.Flags().IntVarP(&fitMinHeight, "min_height", "l", 2, "Min program height used in ramped half-and-half initialization")
	fitCmd.Flags().StringVarP(&fitOutputName, "output", "o", "program.json", "Path where to save the output program")
	fitCmd.Flags().Float64VarP(&fitPLeaf, "p_leaf", "p", 0.5, "Probability of generating a leaf node in ramped half-and-half initialization")
	fitCmd.Flags().Float64VarP(&fitPVariable, "p_variable", "v", 0.5, "Probability of picking a variable and not a constant when generating leaf nodes")
	fitCmd.Flags().IntVarP(&fitRounds, "rounds", "r", 1, "Number of boosting rounds")
	fitCmd.Flags().StringVarP(&fitTargetCol, "target_col", "y", "y", "Name of the target column in the training set")
	fitCmd.Flags().IntVarP(&fitTuningGenerations, "t_generations", "t", 30, "Number of tuning generations")
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
		metric, err := metrics.GetMetric(fitMetricName, fitClass)
		if err != nil {
			return err
		}
		if metric.BiggerIsBetter() {
			metric = metrics.NegativeMetric{Metric: metric}
		}

		// Determine the functions to use
		functions, err := parseStringFuncs(fitFuncsString)
		if err != nil {
			return err
		}

		// Load the training set in memory
		train, err := dataset.ReadCSV(file, fitTargetCol, metric.Classification())
		if err != nil {
			return err
		}

		// Instantiate an Estimator
		estimator := xgp.Estimator{
			Metric:    metric,
			PVariable: fitPVariable,
			NodeInitializer: xgp.RampedHaldAndHalfInitializer{
				MinHeight: fitMinHeight,
				MaxHeight: fitMaxHeight,
				PLeaf:     fitPLeaf,
			},
			Functions:         functions,
			Generations:       fitGenerations,
			TuningGenerations: fitTuningGenerations,
		}

		// Set the initial GA
		estimator.GA = &gago.GA{
			NewGenome: estimator.NewProgram,
			NPops:     1,
			PopSize:   30,
			Model: gago.ModGenerational{
				Selector: gago.SelTournament{
					NContestants: 3,
				},
				MutRate: 0.5,
			},
		}

		err = estimator.Fit(train.X, train.Y, true)

		// // No boosting
		// if fitRounds <= 1 {
		// 	color.Blue("Training without boosting")
		// 	// Monitor progress
		// 	var done = make(chan bool)
		// 	go monitorProgress(estimator.ProgressChan, done)
		// 	// Fit the Estimator
		// 	err = estimator.Fit(train.X, train.Y)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	// Save the best Program
		// 	<-done
		// 	color.Blue("Saving best program to '%s'...", fitOutputName)
		// 	var bestProg, _ = estimator.BestProgram()

		// 	fmt.Println(tree.GetNNodes(bestProg.Root))

		// 	err = xgp.SaveProgramToJSON(*bestProg, fitOutputName)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	return nil
		// }

		return nil
	},
}
