package cmd

import (
	"fmt"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/spf13/cobra"
)

type scoreCmd struct {
	modelPath  string
	targetCol  string
	metricName string

	*cobra.Command
}

func (c *scoreCmd) run(cmd *cobra.Command, args []string) error {
	// Determine the metric to use
	metric, err := metrics.ParseMetric(c.metricName, 1)
	if err != nil {
		return err
	}

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

	// Check the target column exists
	var columns = df.Names()
	if !containsString(columns, c.targetCol) {
		return fmt.Errorf("No column named %s", c.targetCol)
	}

	// Make predictions
	yPred, err := sm.Model.Predict(dataFrameToFloat64(df), metric.NeedsProbabilities())
	if err != nil {
		return err
	}

	// Calculate score
	score, err := metric.Apply(df.Col(c.targetCol).Float(), yPred, nil)
	if err != nil {
		return err
	}
	fmt.Printf("%s: %.5f\n", metric.String(), score)

	return nil
}

func newScoreCmd() *scoreCmd {
	c := &scoreCmd{}
	c.Command = &cobra.Command{
		Use:   "score [dataset path]",
		Short: "Loads a model and computes a score for a given dataset",
		Long:  "Loads a model and computes a score for a given dataset",
		Args:  cobra.ExactArgs(1),
		RunE:  c.run,
	}

	c.Flags().StringVarP(&c.modelPath, "model", "", "model.json", "path to the program to score")
	c.Flags().StringVarP(&c.targetCol, "target", "", "y", "name of the target column in the dataset")
	c.Flags().StringVarP(&c.metricName, "metric", "", "mse", "evaluation metric")

	return c
}
