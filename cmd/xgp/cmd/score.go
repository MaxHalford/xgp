package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/MaxHalford/xgp/meta"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/spf13/cobra"
)

var (
	scoreMetricName   string
	scoreEnsembleName string
	scoreTargetCol    string
)

func init() {
	RootCmd.AddCommand(scoreCmd)

	scoreCmd.Flags().StringVarP(&scoreMetricName, "metric", "", "mse", "evaluation metric")
	scoreCmd.Flags().StringVarP(&scoreEnsembleName, "ensemble", "", "ensemble.json", "path to the program to score")
	scoreCmd.Flags().StringVarP(&scoreTargetCol, "target", "", "y", "name of the target column in the dataset")
}

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Loads an ensemble and computes a score for a given dataset",
	Long:  "Loads an ensemble and computes a score for a given dataset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Load the ensemble
		var (
			ensemble   meta.Ensemble
			bytes, err = ioutil.ReadFile(scoreEnsembleName)
		)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &ensemble)
		if err != nil {
			return err
		}

		// Determine the metric to use
		metric, err := metrics.ParseMetric(scoreMetricName, 1)
		if err != nil {
			return err
		}

		// Load the test set in memory
		df, duration, err := ReadFile(args[0])
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Dataset set took %v to load", duration))

		// Check the target column exists
		var columns = df.Names()
		if !containsString(columns, scoreTargetCol) {
			return fmt.Errorf("No column named %s", scoreTargetCol)
		}

		// Make predictions
		yPred, err := ensemble.Predict(dataFrameToFloat64(df), metric.NeedsProbabilities())
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
