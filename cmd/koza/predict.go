package koza

import (
	"fmt"

	"github.com/MaxHalford/koza"
	"github.com/MaxHalford/koza/dataset"
	"github.com/MaxHalford/koza/metrics"
	"github.com/spf13/cobra"
)

var (
	predictEvalMetricName string
	predictProgramName    string
	predictTargetCol      string
)

func init() {
	RootCmd.AddCommand(predictCmd)

	predictCmd.Flags().StringVarP(&predictEvalMetricName, "eval", "", "mae", "evaluation metric")
	predictCmd.Flags().StringVarP(&predictProgramName, "program", "", "program.json", "path to the program")
	predictCmd.Flags().StringVarP(&predictTargetCol, "target", "", "y", "name of the target column")
}

var predictCmd = &cobra.Command{
	Use:   "predict",
	Short: "Predicts a dataset with a program",
	Long:  "Predicts a dataset with a program",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Determine the metric to use
		metric, err := metrics.GetMetric(predictEvalMetricName, 1)
		if err != nil {
			return err
		}

		// Load the test set in memory
		test, err := dataset.ReadCSV(args[0], predictTargetCol, metric.Classification())
		if err != nil {
			return err
		}

		// Load the program
		prog, err := koza.LoadProgramFromJSON(predictProgramName)
		if err != nil {
			return err
		}

		// Make predictions
		yPred, err := prog.Predict(test.X, metric.NeedsProbabilities())
		if err != nil {
			return err
		}

		fmt.Println(yPred)

		// Calculate score
		score, err := metric.Apply(test.Y, yPred, nil)
		if err != nil {
			return err
		}
		fmt.Printf("Test score: %.5f\n", score)

		return nil
	},
}
