package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MaxHalford/koza"
	"github.com/kniren/gota/series"
	"github.com/spf13/cobra"
)

var (
	predKeptColumns string
	predOutputPath  string
	predProgramName string
	predTargetCol   string
)

func init() {
	RootCmd.AddCommand(predCmd)

	predCmd.Flags().StringVarP(&predKeptColumns, "keep", "", "", "comma-separated columns to keep in the CSV output")
	predCmd.Flags().StringVarP(&predOutputPath, "output", "", "y_pred.csv", "path to the CSV output")
	predCmd.Flags().StringVarP(&predProgramName, "program", "", "program.json", "path to the program used to make predictions")
	predCmd.Flags().StringVarP(&predTargetCol, "target", "", "y", "name of the target column in the CSV output")
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
		df, duration, err := ReadFile(args[0])
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Test set took %v to load", duration))

		// Juggle with the column names
		var (
			featureColumns = df.Names()
			keptColumns    = strings.Split(predKeptColumns, ",")
		)
		for _, col := range keptColumns {
			featureColumns = removeString(featureColumns, col)
		}

		// Make predictions
		yPred, err := prog.Predict(
			dataFrameToFloat64(df.Select(featureColumns)),
			prog.Task.Metric.NeedsProbabilities(),
		)
		if err != nil {
			return err
		}

		// Save the predictions
		outFile, err := os.Create(predOutputPath)
		if err != nil {
			return err
		}
		var pred = df.Select(keptColumns)
		if prog.Task.Metric.Classification() && !prog.Task.Metric.NeedsProbabilities() {
			pred = pred.Mutate(series.New(yPred, series.Int, predTargetCol))
		} else {
			pred = pred.Mutate(series.New(yPred, series.Float, predTargetCol))
		}
		pred.WriteCSV(outFile)

		return nil
	},
}
