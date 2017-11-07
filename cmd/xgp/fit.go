package xgp

import (
	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataset"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/tree"
	"github.com/spf13/cobra"
)

var (
	fitEvalMetricName    string
	fitFuncsString       string
	fitGenerations       int
	fitLossMetricName    string
	fitMaxHeight         int
	fitMinHeight         int
	fitOutputName        string
	fitParsimonyCoeff    float64
	fitPTerminal         float64
	fitPConstant         float64
	fitRounds            int
	fitTargetCol         string
	fitTuningGenerations int
	fitVerbose           bool
)

func init() {
	RootCmd.AddCommand(fitCmd)

	fitCmd.Flags().StringVarP(&fitEvalMetricName, "eval_metric", "e", "", "metric used for monitoring progress, defaults to fit_metric if not provided")
	fitCmd.Flags().StringVarP(&fitLossMetricName, "loss_metric", "l", "mae", "metric used for scoring program, determines the task to perform")
	fitCmd.Flags().StringVarP(&fitFuncsString, "functions", "f", "sum,sub,mul,div,pow,cos,sin", "allowed operators")
	fitCmd.Flags().IntVarP(&fitGenerations, "generations", "g", 30, "number of generations")
	fitCmd.Flags().IntVarP(&fitMaxHeight, "max_height", "b", 6, "max program height used in ramped half-and-half initialization")
	fitCmd.Flags().IntVarP(&fitMinHeight, "min_height", "a", 3, "min program height used in ramped half-and-half initialization")
	fitCmd.Flags().StringVarP(&fitOutputName, "output", "j", "program.json", "path where to save the best program as a JSON file")
	fitCmd.Flags().Float64VarP(&fitParsimonyCoeff, "parsimony", "p", 0, "parsimony coefficient by which a program's height is multiplied to decrease it's fitness")
	fitCmd.Flags().Float64VarP(&fitPTerminal, "p_terminal", "t", 0.5, "probability of generating a terminal branch in ramped half-and-half initialization")
	fitCmd.Flags().Float64VarP(&fitPConstant, "p_constant", "c", 0.5, "probability of picking a constant and not a constant when generating terminal nodes")
	fitCmd.Flags().IntVarP(&fitRounds, "rounds", "r", 1, "number of boosting rounds")
	fitCmd.Flags().StringVarP(&fitTargetCol, "target_col", "y", "y", "name of the target column in the training set")
	fitCmd.Flags().BoolVarP(&fitVerbose, "verbose", "v", true, "monitor progress or not")
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

		// Determine the fitness and evaluation metrics to use
		lossMetric, err := metrics.GetMetric(fitLossMetricName, 1)
		if err != nil {
			return err
		}

		// Default the evaluation metric to the fitness metric if it's nil
		var evalMetric metrics.Metric
		if fitEvalMetricName == "" {
			evalMetric = lossMetric
		} else {
			metric, err := metrics.GetMetric(fitEvalMetricName, 1)
			if err != nil {
				return err
			}
			evalMetric = metric
		}

		// The convention is to use a fitness metric which has to be minimized
		if lossMetric.BiggerIsBetter() {
			lossMetric = metrics.NegativeMetric{Metric: lossMetric}
		}

		// Determine the functions to use
		functions, err := tree.ParseStringFuncs(fitFuncsString)
		if err != nil {
			return err
		}

		// Load the training set in memory
		train, err := dataset.ReadCSV(file, fitTargetCol, lossMetric.Classification())
		if err != nil {
			return err
		}

		// Instantiate an Estimator
		var estimator = xgp.Estimator{
			EvalMetric: evalMetric,
			LossMetric: lossMetric,
			ConstMin:   -10,
			ConstMax:   10,
			PConstant:  fitPConstant,
			TreeInitializer: tree.RampedHaldAndHalfInitializer{
				MinHeight: fitMinHeight,
				MaxHeight: fitMaxHeight,
				PTerminal: fitPTerminal,
			},
			Functions:         functions,
			Generations:       fitGenerations,
			TuningGenerations: fitTuningGenerations,
			ParsimonyCoeff:    fitParsimonyCoeff,
		}

		// Set the initial GA
		estimator.GA = &gago.GA{
			NewGenome: estimator.NewProgram,
			NPops:     1,
			PopSize:   1000,
			Model: gago.ModGenerational{
				Selector: gago.SelTournament{
					NContestants: 3,
				},
				MutRate: 0.5,
			},
		}

		// Fit the estimator
		err = estimator.Fit(train.X, train.Y, train.XNames, fitVerbose)
		if err != nil {
			return err
		}

		// Save the best Program
		bestProg, err := estimator.BestProgram()
		if err != nil {
			return err
		}
		err = xgp.SaveProgramToJSON(*bestProg, fitOutputName)
		if err != nil {
			return err
		}
		return nil
	},
}
