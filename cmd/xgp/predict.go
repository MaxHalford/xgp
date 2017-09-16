package xgp

import (
	"fmt"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/dataframe"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(predictCmd)

	predictCmd.Flags().StringVarP(&programName, "program", "p", "program.json", "Path to the program")
	predictCmd.Flags().StringVarP(&metricName, "metric", "m", "mse", "Metric to use")
	predictCmd.Flags().IntVarP(&class, "class", "c", 1, "Which class to apply the metric to if applicable")
	predictCmd.Flags().StringVarP(&targetCol, "target_col", "y", "y", "Name of the target column")
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
		metric, err := getMetric(metricName, class)
		if err != nil {
			return err
		}

		// Load the test set in memory
		test, err := dataframe.ReadCSV(file, targetCol, false)
		if err != nil {
			return err
		}

		// Load the program
		prog, err := xgp.LoadProgramFromJSON(programName)
		if err != nil {
			return err
		}

		yPred, err := prog.PredictDataFrame(test)
		if err != nil {
			return err
		}
		score, err := metric.Apply(test.Y, yPred, nil)
		if err != nil {
			return err
		}
		fmt.Printf("Test score: %.3f\n", score)

		return nil
	},
}
