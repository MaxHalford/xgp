package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gonum/floats"
	"github.com/kniren/gota/series"
	"github.com/spf13/cobra"
)

type predictCmd struct {
	modelPath  string
	targetCol  string
	proba      bool
	outputPath string
	keepCols   string

	*cobra.Command
}

func (c *predictCmd) run(cmd *cobra.Command, args []string) error {
	// Load the model
	sm, err := readModel(c.modelPath)
	if err != nil {
		return err
	}

	// Load the dataset in memory
	df, duration, err := readFile(args[0])
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Dataset took %v to load", duration))

	// Juggle with the column names
	var (
		featureCols = df.Names()
		keptCols    = strings.Split(c.keepCols, ",")
	)
	for _, col := range keptCols {
		featureCols = removeString(featureCols, col)
	}

	// Check the dataset doesn't contain any missing values
	var XTest = dataFrameToFloat64(df.Select(featureCols))
	for i, x := range XTest {
		if floats.HasNaN(x) {
			return fmt.Errorf("Column '%s' in the test set has missing values", featureCols[i])
		}
	}

	// Make predictions
	yPred, err := sm.Model.Predict(XTest, c.proba)
	if err != nil {
		return err
	}

	// Save the predictions
	outFile, err := os.Create(c.outputPath)
	if err != nil {
		return err
	}
	var pred = df.Select(keptCols).Mutate(series.New(yPred, series.Float, c.targetCol))
	pred.WriteCSV(outFile)

	return nil
}

func newPredictCmd() *predictCmd {
	c := &predictCmd{}
	c.Command = &cobra.Command{
		Use:   "predict [dataset path]",
		Short: "Loads a model and makes predictions for a given dataset",
		Long:  "Loads a model and makes predictions for a given dataset",
		Args:  cobra.ExactArgs(1),
		RunE:  c.run,
	}

	c.Flags().StringVarP(&c.modelPath, "model", "", "model.json", "path to the model used to make predictions")
	c.Flags().StringVarP(&c.targetCol, "target", "", "y", "name of the target column in the CSV output")
	c.Flags().BoolVarP(&c.proba, "proba", "", false, "predict probabilities in case of classification")
	c.Flags().StringVarP(&c.outputPath, "output", "", "y_pred.csv", "path to the CSV output")
	c.Flags().StringVarP(&c.keepCols, "keep", "", "", "comma-separated columns to keep in the CSV output")

	return c
}
