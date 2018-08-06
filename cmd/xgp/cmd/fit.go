package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/meta"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/gonum/floats"
	"github.com/spf13/cobra"
)

var (
	// Learning parameters
	fitLossMetricName string
	fitEvalMetricName string
	fitParsimonyCoeff float64
	fitPolishBest     bool

	// Function parameters
	fitFuncs     string
	fitConstMin  float64
	fitConstMax  float64
	fitPConst    float64
	fitPFull     float64
	fitPLeaf     float64
	fitMinHeight uint
	fitMaxHeight uint

	// Genetic algorithm parameters
	fitNPopulations      uint
	fitNIndividuals      uint
	fitNGenerations      uint
	fitPHoistMutation    float64
	fitPSubtreeMutation  float64
	fitPPointMutation    float64
	fitPointMutationRate float64
	fitPSubtreeCrossover float64

	// Other
	fitSeed int64

	// Ensemble parameters
	fitMode                string
	fitLearningRate        float64
	fitNPrograms           uint
	fitRowSampling         float64
	fitColSampling         float64
	fitEarlyStoppingRounds uint

	// CLI parameters
	fitIgnoredCols string
	fitOutputName  string
	fitTargetCol   string
	fitValPath     string
	fitVerbose     bool
)

func init() {
	RootCmd.AddCommand(fitCmd)

	fitCmd.Flags().StringVarP(&fitLossMetricName, "loss", "", "mae", "metric used for scoring program; determines the task to perform")
	fitCmd.Flags().StringVarP(&fitEvalMetricName, "eval", "", "", "metric used for monitoring progress; defaults to loss_metric if not provided")
	fitCmd.Flags().Float64VarP(&fitParsimonyCoeff, "parsimony", "", 0.00001, "parsimony coefficient by which a program's height is multiplied to decrease it's fitness")
	fitCmd.Flags().BoolVarP(&fitPolishBest, "polish", "", true, "whether or not to polish the best program")

	fitCmd.Flags().StringVarP(&fitFuncs, "funcs", "", "add,sub,mul,div", "comma-separated set of authorised functions")
	fitCmd.Flags().Float64VarP(&fitConstMin, "const_min", "", -5, "lower bound used for generating random constants")
	fitCmd.Flags().Float64VarP(&fitConstMax, "const_max", "", 5, "upper bound used for generating random constants")
	fitCmd.Flags().Float64VarP(&fitPConst, "p_const", "", 0.5, "probability of generating a constant instead of a variable")
	fitCmd.Flags().Float64VarP(&fitPFull, "p_full", "", 0.5, "probability of using full initialization during ramped half-and-half initialization")
	fitCmd.Flags().Float64VarP(&fitPLeaf, "p_leaf", "", 0.3, "probability of generating a terminal node during ramped half-and-half initialization")
	fitCmd.Flags().UintVarP(&fitMinHeight, "min_height", "", 3, "minimum program height used in ramped half-and-half initialization")
	fitCmd.Flags().UintVarP(&fitMaxHeight, "max_height", "", 5, "maximum program height used in ramped half-and-half initialization")

	fitCmd.Flags().UintVarP(&fitNPopulations, "pops", "", 1, "number of populations used in the GA")
	fitCmd.Flags().UintVarP(&fitNIndividuals, "indis", "", 100, "number of individuals used for each population in the GA")
	fitCmd.Flags().UintVarP(&fitNGenerations, "gens", "", 30, "number of generations used in the GA")
	fitCmd.Flags().Float64VarP(&fitPHoistMutation, "p_hoist_mut", "", 0.1, "probability of applying hoist mutation")
	fitCmd.Flags().Float64VarP(&fitPSubtreeMutation, "p_sub_mut", "", 0.1, "probability of applying subtree mutation")
	fitCmd.Flags().Float64VarP(&fitPPointMutation, "p_point_mut", "", 0.1, "probability of applying point mutation")
	fitCmd.Flags().Float64VarP(&fitPointMutationRate, "point_mut_rate", "", 0.3, "probability of modifying an operator during point mutation")
	fitCmd.Flags().Float64VarP(&fitPSubtreeCrossover, "p_sub_cross", "", 0.5, "probability of applying subtree crossover")

	fitCmd.Flags().Int64VarP(&fitSeed, "seed", "", 0, "seed for random number generation")

	fitCmd.Flags().StringVarP(&fitMode, "mode", "", "adaboost", "training mode")
	fitCmd.Flags().UintVarP(&fitNPrograms, "n_programs", "", 30, "number of programs to use in the ensemble")
	fitCmd.Flags().Float64VarP(&fitLearningRate, "learning_rate", "", 0.05, "learning rate used for boosting")
	fitCmd.Flags().Float64VarP(&fitRowSampling, "row_sampling", "", 0.6, "row sampling used for bagging")
	fitCmd.Flags().Float64VarP(&fitColSampling, "col_sampling", "", 1, "column sampling used for bagging")
	fitCmd.Flags().UintVarP(&fitEarlyStoppingRounds, "early_stopping", "", 5, "number of rounds after which training stops if the evaluation score worsens")

	fitCmd.Flags().StringVarP(&fitIgnoredCols, "ignore", "", "", "comma-separated columns to ignore")
	fitCmd.Flags().StringVarP(&fitOutputName, "output", "", "ensemble.json", "path where to save the JSON representation of the best program ")
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

		// Determine the loss metric
		lossMetric, err := metrics.ParseMetric(fitLossMetricName, 1)
		if err != nil {
			return err
		}

		// Determine the evaluation metric
		if fitEvalMetricName == "" {
			fitEvalMetricName = fitLossMetricName
		}
		evalMetric, err := metrics.ParseMetric(fitEvalMetricName, 1)
		if err != nil {
			return err
		}

		// Instantiate an GP
		var config = xgp.GPConfig{
			LossMetric:     lossMetric,
			EvalMetric:     evalMetric,
			ParsimonyCoeff: fitParsimonyCoeff,
			PolishBest:     fitPolishBest,

			Funcs:     fitFuncs,
			ConstMin:  fitConstMin,
			ConstMax:  fitConstMax,
			PConst:    fitPConst,
			PFull:     fitPFull,
			PLeaf:     fitPLeaf,
			MinHeight: fitMinHeight,
			MaxHeight: fitMaxHeight,

			NPopulations:      fitNPopulations,
			NIndividuals:      fitNIndividuals,
			NGenerations:      fitNGenerations,
			PHoistMutation:    fitPHoistMutation,
			PSubtreeMutation:  fitPSubtreeMutation,
			PPointMutation:    fitPPointMutation,
			PointMutationRate: fitPointMutationRate,
			PSubtreeCrossover: fitPSubtreeCrossover,

			RNG: rng,
		}

		// Load the training set in memory
		train, duration, err := ReadFile(args[0])
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Training dataset took %v to load", duration))

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

		// Check XTrain doesn't contain any NaNs
		for i, x := range XTrain {
			if floats.HasNaN(x) {
				return fmt.Errorf("Column '%s' in the training set has NaNs", featureColumns[i])
			}
		}

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
			fmt.Println(fmt.Sprintf("Validation dataset took %v to load", duration))
			XVal = dataFrameToFloat64(val.Select(featureColumns))
			YVal = val.Col(fitTargetCol).Float()
		}

		// Check XVal doesn't contain any NaNs
		for i, x := range XVal {
			if floats.HasNaN(x) {
				return fmt.Errorf("Column '%s' in the validation set has NaNs", featureColumns[i])
			}
		}

		// Train
		gb, err := meta.NewGradBoost(
			config,
			fitNPrograms,
			fitEarlyStoppingRounds,
			fitLearningRate,
			meta.SquareLoss{},
		)
		if err != nil {
			return err
		}
		err = gb.Fit(XTrain, YTrain, XVal, YVal)
		if err != nil {
			return err
		}

		// Save the ensemble
		/*bytes, err := json.Marshal(ensemble)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(fitOutputName, bytes, 0644)
		if err != nil {
			return err
		}*/

		return nil
	},
}
