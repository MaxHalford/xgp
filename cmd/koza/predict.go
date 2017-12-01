package koza

import (
	"fmt"

	"github.com/MaxHalford/koza"
	"github.com/spf13/cobra"
)

var (
	predictOutputName  string
	predictProgramName string
)

func init() {
	RootCmd.AddCommand(predictCmd)

	predictCmd.Flags().StringVarP(&predictOutputName, "output", "", "y_pred.csv", "path to the CSV output")
	predictCmd.Flags().StringVarP(&predictProgramName, "program", "", "program.json", "path to the program")
}

var predictCmd = &cobra.Command{
	Use:   "predict",
	Short: "Predicts a dataset with a program",
	Long:  "Predicts a dataset with a program",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Load the program
		prog, err := koza.LoadProgramFromJSON(predictProgramName)
		if err != nil {
			return err
		}

		// Load the test set in memory
		df, err := ReadCSV(args[0])
		if err != nil {
			return err
		}

		// Make predictions
		yPred, err := prog.Predict(dataFrameToFloat64(df), prog.Task.Metric.NeedsProbabilities())
		if err != nil {
			return err
		}

		// Save predictions
		fmt.Println(yPred)

		return nil
	},
}
