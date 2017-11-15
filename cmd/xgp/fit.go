package xgp

import (
	"github.com/spf13/cobra"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataset"
)

var (
	fitConstMax          float64
	fitConstMin          float64
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
	fitSeed              int64
	fitTargetCol         string
	fitTuningGenerations int
	fitVerbose           bool
)

func init() {
	RootCmd.AddCommand(fitCmd)

	fitCmd.Flags().Float64VarP(&fitConstMax, "const_max", "", 5, "upper bound for generating random constants")
	fitCmd.Flags().Float64VarP(&fitConstMin, "const_min", "", -5, "lower bound for generating random constants")
	fitCmd.Flags().StringVarP(&fitEvalMetricName, "eval_metric", "", "", "metric used for monitoring progress, defaults to fit_metric if not provided")
	fitCmd.Flags().StringVarP(&fitLossMetricName, "loss_metric", "", "mae", "metric used for scoring program, determines the task to perform")
	fitCmd.Flags().StringVarP(&fitFuncsString, "functions", "", "sum,sub,mul,div", "allowed operators")
	fitCmd.Flags().IntVarP(&fitGenerations, "generations", "", 30, "number of generations")
	fitCmd.Flags().IntVarP(&fitMaxHeight, "max_height", "", 6, "max program height used in ramped half-and-half initialization")
	fitCmd.Flags().IntVarP(&fitMinHeight, "min_height", "", 3, "min program height used in ramped half-and-half initialization")
	fitCmd.Flags().StringVarP(&fitOutputName, "output", "", "program.json", "path where to save the best program as a JSON file")
	fitCmd.Flags().Float64VarP(&fitParsimonyCoeff, "parsimony", "", 0, "parsimony coefficient by which a program's height is multiplied to decrease it's fitness")
	fitCmd.Flags().Float64VarP(&fitPTerminal, "p_terminal", "", 0.5, "probability of generating a terminal branch in ramped half-and-half initialization")
	fitCmd.Flags().Float64VarP(&fitPConstant, "p_constant", "", 0.5, "probability of picking a constant and not a constant when generating terminal nodes")
	fitCmd.Flags().IntVarP(&fitRounds, "rounds", "", 1, "number of boosting rounds")
	fitCmd.Flags().Int64VarP(&fitSeed, "seed", "", 0, "seed for random number generation")
	fitCmd.Flags().StringVarP(&fitTargetCol, "target_col", "", "y", "name of the target column in the training set")
	fitCmd.Flags().BoolVarP(&fitVerbose, "verbose", "", true, "monitor progress or not")
}

var fitCmd = &cobra.Command{
	Use:   "fit",
	Short: "Fits a program to a dataset",
	Long:  "Fits a program to a dataset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Instantiate an Estimator
		var estimator, err = xgp.NewEstimator(
			fitConstMax,
			fitConstMin,
			fitEvalMetricName,
			fitFuncsString,
			fitLossMetricName,
			fitMaxHeight,
			fitMinHeight,
			fitGenerations,
			fitParsimonyCoeff,
			fitPConstant,
			fitPTerminal,
			fitSeed,
			fitTuningGenerations,
		)
		if err != nil {
			return err
		}

		// Load the training set in memory
		train, err := dataset.ReadCSV(args[0], fitTargetCol, estimator.LossMetric.Classification())
		if err != nil {
			return err
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
