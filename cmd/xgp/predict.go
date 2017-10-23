package xgp

import (
	"fmt"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataset"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	predictClass       int
	predictMetricName  string
	predictProgramName string
	predictTargetCol   string
)

func init() {
	RootCmd.AddCommand(predictCmd)

	predictCmd.Flags().IntVarP(&predictClass, "class", "c", 1, "Which class to apply the metric to if applicable")
	predictCmd.Flags().StringVarP(&predictMetricName, "metric", "m", "mse", "Metric to use")
	predictCmd.Flags().StringVarP(&predictProgramName, "program", "p", "program.json", "Path to the program")
	predictCmd.Flags().StringVarP(&predictTargetCol, "target_col", "y", "y", "Name of the target column")
}

var predictCmd = &cobra.Command{
	Use:   "predict",
	Short: "Predicts a dataset with a program",
	Long:  "Predicts a dataset with a program",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if the file exists
		var file = args[0]
		if err := fileExists(file); err != nil {
			return err
		}

		// Determine the metric to use
		metric, err := metrics.GetMetric(predictMetricName, predictClass)
		if err != nil {
			return err
		}

		// Load the test set in memory
		test, err := dataset.ReadCSV(file, predictTargetCol, metric.Classification())
		if err != nil {
			return err
		}

		// Load the program
		prog, err := xgp.LoadProgramFromJSON(predictProgramName)
		if err != nil {
			return err
		}

		yPred, err := prog.Predict(test.X)
		if err != nil {
			return err
		}
		fmt.Println(yPred)
		score, err := metric.Apply(test.Y, yPred, nil)
		if err != nil {
			return err
		}
		color.Green("Test score: %.5f\n", score)

		return nil
	},
}
