package koza

import (
	"os"

	"github.com/MaxHalford/koza"
	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"github.com/spf13/cobra"
)

var (
	predOutputPath  string
	predProgramName string
	predTargetCol   string
)

func init() {
	RootCmd.AddCommand(predCmd)

	predCmd.Flags().StringVarP(&predOutputPath, "output", "", "y_pred.csv", "path to the CSV output")
	predCmd.Flags().StringVarP(&predProgramName, "program", "", "program.json", "path to the program")
	predCmd.Flags().StringVarP(&predTargetCol, "target", "", "y", "name of the predictions column in the CSV output")
}

var predCmd = &cobra.Command{
	Use:   "predict",
	Short: "Predicts a dataset with a program",
	Long:  "Predicts a dataset with a program",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Load the program
		prog, err := koza.LoadProgramFromJSON(predProgramName)
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

		// Save the predictions
		outFile, err := os.Create(predOutputPath)
		if err != nil {
			return err
		}
		df = dataframe.New(series.New(yPred, series.Float, predTargetCol))
		df.WriteCSV(outFile)

		return nil
	},
}
