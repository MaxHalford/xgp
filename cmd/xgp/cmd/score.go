package cmd

import (
	"fmt"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/spf13/cobra"
)

var (
	scoreEvalMetricName string
	scoreProgramName    string
	scoreTargetCol      string
)

func init() {
	RootCmd.AddCommand(scoreCmd)

	scoreCmd.Flags().StringVarP(&scoreEvalMetricName, "eval", "", "", "evaluation metric")
	scoreCmd.Flags().StringVarP(&scoreProgramName, "program", "", "program.json", "path to the program to score")
	scoreCmd.Flags().StringVarP(&scoreTargetCol, "target", "", "y", "name of the target column in the dataset")
}

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Predicts a dataset with a program",
	Long:  "Predicts a dataset with a program",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Load the program
		prog, err := xgp.LoadProgramFromJSON(scoreProgramName)
		if err != nil {
			return err
		}

		// Determine the metric to use
		var metric metrics.Metric
		if scoreEvalMetricName == "" {
			metric, err = metrics.GetMetric(scoreEvalMetricName, 1)
			if err != nil {
				return err
			}
		} else {
			metric = prog.Task.Metric
		}

		// Load the test set in memory
		df, duration, err := ReadFile(args[0])
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Scoring set took %v to load", duration))

		// Check the target column exists
		var columns = df.Names()
		if !containsString(columns, scoreTargetCol) {
			return fmt.Errorf("No column named %s", scoreTargetCol)
		}

		// Make predictions
		yPred, err := prog.Predict(dataFrameToFloat64(df), metric.NeedsProbabilities())
		if err != nil {
			return err
		}

		// Calculate score
		score, err := metric.Apply(df.Col(scoreTargetCol).Float(), yPred, nil)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %.5f\n", metric.String(), score)

		return nil
	},
}
