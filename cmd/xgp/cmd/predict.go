package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/MaxHalford/xgp/meta"
	"github.com/gonum/floats"
	"github.com/kniren/gota/series"
	"github.com/spf13/cobra"
)

var (
	predPredictProba bool
	predKeptColumns  string
	predOutputPath   string
	predEnsembleName string
	predTargetCol    string
)

func init() {
	RootCmd.AddCommand(predCmd)

	predCmd.Flags().BoolVarP(&predPredictProba, "proba", "", false, "predict probabilities in case of classification")
	predCmd.Flags().StringVarP(&predKeptColumns, "keep", "", "", "comma-separated columns to keep in the CSV output")
	predCmd.Flags().StringVarP(&predOutputPath, "output", "", "y_pred.csv", "path to the CSV output")
	predCmd.Flags().StringVarP(&predEnsembleName, "ensemble", "", "ensemble.json", "path to the model used to make predictions")
	predCmd.Flags().StringVarP(&predTargetCol, "target", "", "y", "name of the target column in the CSV output")
}

var predCmd = &cobra.Command{
	Use:   "predict",
	Short: "Loads an ensemble and makes predictions on a given dataset",
	Long:  "Loads an ensemble and makes predictions on a given dataset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Load the ensemble
		var (
			ensemble   meta.Ensemble
			bytes, err = ioutil.ReadFile(predEnsembleName)
		)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &ensemble)
		if err != nil {
			return err
		}

		// Load the test set in memory
		df, duration, err := ReadFile(args[0])
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Dataset set took %v to load", duration))

		// Juggle with the column names
		var (
			featureColumns = df.Names()
			keptColumns    = strings.Split(predKeptColumns, ",")
		)
		for _, col := range keptColumns {
			featureColumns = removeString(featureColumns, col)
		}

		// Check XVal doesn't contain any NaNs
		var XTest = dataFrameToFloat64(df.Select(featureColumns))
		for i, x := range XTest {
			if floats.HasNaN(x) {
				return fmt.Errorf("Column '%s' in the test set has NaNs", featureColumns[i])
			}
		}

		// Make predictions
		yPred, err := ensemble.Predict(XTest, predPredictProba)
		if err != nil {
			return err
		}

		// Save the predictions
		outFile, err := os.Create(predOutputPath)
		if err != nil {
			return err
		}
		var pred = df.Select(keptColumns)
		if ensemble.Classification && !predPredictProba {
			pred = pred.Mutate(series.New(yPred, series.Int, predTargetCol))
		} else {
			pred = pred.Mutate(series.New(yPred, series.Float, predTargetCol))
		}
		pred.WriteCSV(outFile)

		return nil
	},
}
