package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/ensemble"
	"github.com/spf13/cobra"
)

var (
	fitConstMin float64
	fitConstMax float64

	fitEvalMetricName string
	fitLossMetricName string

	fitFuncs string

	fitMinHeight int
	fitMaxHeight int

	fitNPopulations       int
	fitNIndividuals       int
	fitNGenerations       int
	fitNTuningGenerations int
	fitNRounds            int

	fitPConstant float64
	fitPFull     float64
	fitPTerminal float64

	fitPHoistMutation    float64
	fitPPointMutation    float64
	fitPSubTreeMutation  float64
	fitPointMutationRate float64

	fitPSubTreeCrossover float64
	fitParsimonyCoeff    float64

	fitSeed int64

	fitIgnoredCols string
	fitOutputName  string
	fitTargetCol   string
	fitValPath     string
	fitVerbose     bool
)

func init() {
	RootCmd.AddCommand(fitCmd)

	fitCmd.Flags().Float64VarP(&fitConstMin, "const_min", "", -5, "lower bound for generating random constants")
	fitCmd.Flags().Float64VarP(&fitConstMax, "const_max", "", 5, "upper bound for generating random constants")

	fitCmd.Flags().StringVarP(&fitEvalMetricName, "eval", "", "", "metric used for monitoring progress, defaults to fit_metric if not provided")
	fitCmd.Flags().StringVarP(&fitLossMetricName, "loss", "", "mae", "metric used for scoring program, determines the task to perform")

	fitCmd.Flags().StringVarP(&fitFuncs, "funcs", "", "sum,sub,mul,div", "comma-separated authorised operators")

	fitCmd.Flags().IntVarP(&fitMinHeight, "min_height", "", 3, "min program height used in ramped half-and-half initialization")
	fitCmd.Flags().IntVarP(&fitMaxHeight, "max_height", "", 5, "max program height used in ramped half-and-half initialization")

	fitCmd.Flags().IntVarP(&fitNPopulations, "pops", "", 1, "number of populations to use in the GA")
	fitCmd.Flags().IntVarP(&fitNIndividuals, "indis", "", 50, "number of individuals to use for each population in the GA")
	fitCmd.Flags().IntVarP(&fitNGenerations, "gens", "", 30, "number of generations used to evolve the GA")
	fitCmd.Flags().IntVarP(&fitNTuningGenerations, "tune_gens", "", 0, "number of generations used to evolve the tuning GA")
	fitCmd.Flags().IntVarP(&fitNRounds, "rounds", "", 1, "number of programs used in the ensemble")

	fitCmd.Flags().Float64VarP(&fitPConstant, "p_constant", "", 0.5, "probability of picking a constant and not a constant when generating terminal nodes")
	fitCmd.Flags().Float64VarP(&fitPFull, "p_full", "", 0.5, "probability of use full initialization during ramped half-and-half initialization")
	fitCmd.Flags().Float64VarP(&fitPTerminal, "p_terminal", "", 0.3, "probability of generating a terminal branch in ramped half-and-half initialization")

	fitCmd.Flags().Float64VarP(&fitPHoistMutation, "p_hoist_mut", "", 0.1, "probability of applying hoist mutation")
	fitCmd.Flags().Float64VarP(&fitPSubTreeMutation, "p_sub_mut", "", 0.1, "probability of applying sub-tree mutation")
	fitCmd.Flags().Float64VarP(&fitPPointMutation, "p_point_mut", "", 0.1, "probability of applying point mutation")
	fitCmd.Flags().Float64VarP(&fitPointMutationRate, "point_mut_rate", "", 0.3, "probability of modifying an operator during point mutation")

	fitCmd.Flags().Float64VarP(&fitPSubTreeCrossover, "p_sub_cross", "", 0.5, "probability of applying sub-tree crossover")

	fitCmd.Flags().Float64VarP(&fitParsimonyCoeff, "parsimony", "", 0, "parsimony coefficient by which a program's height is multiplied to decrease it's fitness")

	fitCmd.Flags().Int64VarP(&fitSeed, "seed", "", 0, "seed for random number generation")

	fitCmd.Flags().StringVarP(&fitIgnoredCols, "ignore", "", "", "comma-separated columns to ignore")
	fitCmd.Flags().StringVarP(&fitOutputName, "output", "", "program.json", "path where to save the JSON representation of the best program ")
	fitCmd.Flags().StringVarP(&fitTargetCol, "target", "", "y", "name of the target column in the training and validation datasets")
	fitCmd.Flags().StringVarP(&fitValPath, "val", "", "", "path to a validation dataset that can be used to monitor out-of-bag performance")
	fitCmd.Flags().BoolVarP(&fitVerbose, "verbose", "", true, "display progress in the shell")
}

var fitCmd = &cobra.Command{
	Use:   "fit",
	Short: "Fits a program to a dataset",
	Long:  "Fits a program to a dataset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Instantiate a random number generator
		var rng *rand.Rand
		if fitSeed == 0 {
			rng = rand.New(rand.NewSource(time.Now().UnixNano()))
		} else {
			rng = rand.New(rand.NewSource(fitSeed))
		}

		// By default the evaluation metric is equal to the loss metric
		if fitEvalMetricName == "" {
			fitEvalMetricName = fitLossMetricName
		}

		// Instantiate an Estimator
		var config = xgp.Config{
			ConstMin: fitConstMin,
			ConstMax: fitConstMax,

			EvalMetricName: fitEvalMetricName,
			LossMetricName: fitLossMetricName,

			Funcs: fitFuncs,

			MinHeight: fitMinHeight,
			MaxHeight: fitMaxHeight,

			NPopulations:       fitNPopulations,
			NIndividuals:       fitNIndividuals,
			NGenerations:       fitNGenerations,
			NTuningGenerations: fitNTuningGenerations,

			PConstant: fitPConstant,
			PFull:     fitPFull,
			PTerminal: fitPTerminal,

			PHoistMutation:    fitPHoistMutation,
			PPointMutation:    fitPPointMutation,
			PSubTreeMutation:  fitPSubTreeMutation,
			PointMutationRate: fitPointMutationRate,

			PSubTreeCrossover: fitPSubTreeCrossover,

			ParsimonyCoeff: fitParsimonyCoeff,

			RNG: rng,
		}
		var estimator, err = config.NewEstimator()
		if err != nil {
			return err
		}

		// Load the training set in memory
		train, duration, err := ReadFile(args[0])
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Training set took %v to load", duration))

		// Check the target column exists
		var columns = train.Names()
		if !containsString(columns, fitTargetCol) {
			return fmt.Errorf("No column named %s", fitTargetCol)
		}
		var featureColumns = removeString(columns, fitTargetCol)

		// Remove ignored columns
		for _, col := range strings.Split(fitIgnoredCols, ",") {
			featureColumns = removeString(featureColumns, col)
		}

		// Extract the features and target from the training set
		var (
			XTrain = dataFrameToFloat64(train.Select(featureColumns))
			YTrain = train.Col(fitTargetCol).Float()
		)

		// Load the validation set in memory
		var (
			XVal [][]float64
			YVal []float64
		)
		if fitValPath != "" {
			val, duration, err := ReadFile(fitValPath)
			if err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("Validation set took %v to load", duration))
			XVal = dataFrameToFloat64(val.Select(featureColumns))
			YVal = val.Col(fitTargetCol).Float()
		}

		// Instantiate an ensemble
		var bag = ensemble.BaggingRegressor{
			NEstimators: fitNRounds,
			RowSampling: 1,
			ColSampling: 1,
			RNG:         rng,
		}

		// Fit the ensemble
		err = bag.Fit(
			estimator,
			XTrain,
			YTrain,
			nil,
			XVal,
			YVal,
			nil,
			fitVerbose,
			estimator.EvalMetric,
			rng,
		)
		if err != nil {
			return err
		}

		// Save the model
		err = ensemble.SaveEnsembleToJSON(&bag, fitOutputName)
		if err != nil {
			return err
		}
		return nil
	},
}
