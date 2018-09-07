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

type fitCmd struct {
	flavor string

	// GP learning parameters
	lossName       string
	evalName       string
	parsimonyCoeff float64
	polishBest     bool
	funcs          string
	constMin       float64
	constMax       float64
	pConst         float64
	pFull          float64
	pLeaf          float64
	minHeight      uint
	maxHeight      uint

	// GA parameters
	nPops         uint
	popSize       uint
	nGenerations  uint
	pHoistMut     float64
	pSubtreeMut   float64
	pPointMut     float64
	pointMutRate  float64
	pSubtreeCross float64

	// Ensemble learning parameters
	nRounds              uint
	nEarlyStoppingRounds uint
	learningRate         float64
	lineSearch           bool

	// Other
	seed int64

	// CLI parameters
	ignoredCols string
	outputPath  string
	targetCol   string
	valPath     string
	verbose     bool

	*cobra.Command
}

func (c *fitCmd) run(cmd *cobra.Command, args []string) error {
	// Instantiate a random number generator
	var rng *rand.Rand
	if c.seed == 0 {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		rng = rand.New(rand.NewSource(c.seed))
	}

	// Determine the loss metric
	lossMetric, err := metrics.ParseMetric(c.lossName, 1)
	if err != nil {
		return err
	}

	// Determine the evaluation metric
	if c.evalName == "" {
		c.evalName = c.lossName
	}
	evalMetric, err := metrics.ParseMetric(c.evalName, 1)
	if err != nil {
		return err
	}

	// Instantiate a GP
	var config = xgp.GPConfig{
		LossMetric:     lossMetric,
		EvalMetric:     evalMetric,
		ParsimonyCoeff: c.parsimonyCoeff,
		PolishBest:     c.polishBest,

		Funcs:     c.funcs,
		ConstMin:  c.constMin,
		ConstMax:  c.constMax,
		PConst:    c.pConst,
		PFull:     c.pFull,
		PLeaf:     c.pLeaf,
		MinHeight: c.minHeight,
		MaxHeight: c.maxHeight,

		NPopulations:      c.nPops,
		NIndividuals:      c.popSize,
		NGenerations:      c.nGenerations,
		PHoistMutation:    c.pHoistMut,
		PSubtreeMutation:  c.pSubtreeMut,
		PPointMutation:    c.pPointMut,
		PointMutationRate: c.pointMutRate,
		PSubtreeCrossover: c.pSubtreeCross,

		RNG: rng,
	}

	// Load the training set in memory
	train, duration, err := readFile(args[0])
	if err != nil {
		return err
	}
	if c.verbose {
		fmt.Printf("Training dataset took %v to load\n", duration)
	}

	// Check the target column exists
	var columns = train.Names()
	if !containsString(columns, c.targetCol) {
		return fmt.Errorf("No column named %s", c.targetCol)
	}
	var featureCols = removeString(columns, c.targetCol)

	// Remove ignored columns
	for _, col := range strings.Split(c.ignoredCols, ",") {
		featureCols = removeString(featureCols, col)
	}

	// Extract the features and target from the training set
	var (
		XTrain = dataFrameToFloat64(train.Select(featureCols))
		YTrain = train.Col(c.targetCol).Float()
	)

	// Check XTrain doesn't contain any NaNs
	for i, x := range XTrain {
		if floats.HasNaN(x) {
			return fmt.Errorf("Column '%s' in the training set has NaNs", featureCols[i])
		}
	}

	// Load the validation set in memory
	var (
		XVal [][]float64
		YVal []float64
	)
	if c.valPath != "" {
		val, duration, err := readFile(c.valPath)
		if err != nil {
			return err
		}
		if c.verbose {
			fmt.Printf("Validation dataset took %v to load\n", duration)
		}
		XVal = dataFrameToFloat64(val.Select(featureCols))
		YVal = val.Col(c.targetCol).Float()
	}

	// Check XVal doesn't contain any NaNs
	for i, x := range XVal {
		if floats.HasNaN(x) {
			return fmt.Errorf("Column '%s' in the validation set has NaNs", featureCols[i])
		}
	}

	// Train
	switch c.flavor {

	case "vanilla":
		gp, err := config.NewGP()
		if err != nil {
			return err
		}
		err = gp.Fit(XTrain, YTrain, nil, XVal, YVal, nil, c.verbose)
		if err != nil {
			return err
		}
		best, err := gp.BestProgram()
		if err != nil {
			return err
		}
		return writeProgram(best, c.outputPath)

	case "boosting":
		loss, ok := lossMetric.(metrics.DiffMetric)
		if !ok {
			return fmt.Errorf("The '%s' metric can't be used for gradient boosting because it is"+
				" not differentiable", lossMetric.String())
		}
		var ls meta.LineSearcher
		if c.lineSearch {
			ls = meta.GoldenLineSearch{
				Min: 0,
				Max: 10,
				Tol: 1e-10,
			}
		}
		gb, err := meta.NewGradientBoosting(
			config,
			c.nRounds,
			c.nEarlyStoppingRounds,
			c.learningRate,
			ls,
			loss,
		)
		if err != nil {
			return err
		}
		err = gb.Fit(XTrain, YTrain, nil, XVal, YVal, nil, c.verbose)
		if err != nil {
			return err
		}
		return writeGradientBoosting(gb, c.outputPath)

	}

	return errUnknownFlavor{c.flavor}
}

func newFitCmd() *fitCmd {
	c := &fitCmd{}
	c.Command = &cobra.Command{
		Use:   "fit [dataset path]",
		Short: "Fits an ensemble of programs to a dataset",
		Long:  "Fits an ensemble of programs to a dataset",
		Args:  cobra.ExactArgs(1),
		RunE:  c.run,
	}

	c.Flags().StringVarP(&c.flavor, "flavor", "", "boosting", "training flavor to use ('boosting' or 'vanilla')")

	c.Flags().StringVarP(&c.lossName, "loss", "", "mse", "metric used for scoring program; determines the task to perform")
	c.Flags().StringVarP(&c.evalName, "eval", "", "", "metric used for monitoring progress; defaults to loss_metric if not provided")
	c.Flags().Float64VarP(&c.parsimonyCoeff, "parsimony", "", 0.00001, "parsimony coefficient by which a program's height is multiplied to decrease it's fitness")
	c.Flags().BoolVarP(&c.polishBest, "polish", "", true, "whether or not to polish the best program")
	c.Flags().StringVarP(&c.funcs, "funcs", "", "min,max,add,sub,mul,div", "comma-separated set of authorised functions")
	c.Flags().Float64VarP(&c.constMin, "const_min", "", -5, "lower bound used for generating random constants")
	c.Flags().Float64VarP(&c.constMax, "const_max", "", 5, "upper bound used for generating random constants")
	c.Flags().Float64VarP(&c.pConst, "p_const", "", 0.5, "probability of generating a constant instead of a variable")
	c.Flags().Float64VarP(&c.pFull, "p_full", "", 0.5, "probability of using full initialization during ramped half-and-half initialization")
	c.Flags().Float64VarP(&c.pLeaf, "p_leaf", "", 0.3, "probability of generating a terminal node during ramped half-and-half initialization")
	c.Flags().UintVarP(&c.minHeight, "min_height", "", 3, "minimum program height used in ramped half-and-half initialization")
	c.Flags().UintVarP(&c.maxHeight, "max_height", "", 5, "maximum program height used in ramped half-and-half initialization")

	c.Flags().UintVarP(&c.nPops, "pops", "", 1, "number of populations used in the GA")
	c.Flags().UintVarP(&c.popSize, "indis", "", 50, "number of individuals used for each population in the GA")
	c.Flags().UintVarP(&c.nGenerations, "gens", "", 30, "number of generations used in the GA")
	c.Flags().Float64VarP(&c.pHoistMut, "p_hoist_mut", "", 0.1, "probability of applying hoist mutation")
	c.Flags().Float64VarP(&c.pSubtreeMut, "p_sub_mut", "", 0.1, "probability of applying subtree mutation")
	c.Flags().Float64VarP(&c.pPointMut, "p_point_mut", "", 0.1, "probability of applying point mutation")
	c.Flags().Float64VarP(&c.pointMutRate, "point_mut_rate", "", 0.3, "probability of modifying an operator during point mutation")
	c.Flags().Float64VarP(&c.pSubtreeCross, "p_sub_cross", "", 0.5, "probability of applying subtree crossover")

	c.Flags().UintVarP(&c.nRounds, "rounds", "", 50, "number of programs to use in case of using an ensemble")
	c.Flags().UintVarP(&c.nEarlyStoppingRounds, "early_stopping", "", 5, "number of rounds after which training stops if the evaluation score worsens")
	c.Flags().Float64VarP(&c.learningRate, "learning_rate", "", 0.08, "learning rate used for boosting")
	c.Flags().BoolVarP(&c.lineSearch, "line_search", "", true, "whether to use line-search or not")

	c.Flags().Int64VarP(&c.seed, "seed", "", 0, "seed for random number generation")

	c.Flags().StringVarP(&c.ignoredCols, "ignore", "", "", "comma-separated columns to ignore")
	c.Flags().StringVarP(&c.outputPath, "output", "", "model.json", "path where to save the JSON representation of the final model")
	c.Flags().StringVarP(&c.targetCol, "target", "", "y", "name of the target column in the training and validation datasets")
	c.Flags().StringVarP(&c.valPath, "val", "", "", "path to a validation dataset that can be used to monitor out-of-bag performance")
	c.Flags().BoolVarP(&c.verbose, "verbose", "", true, "whether to display progress in the shell or not")

	return c
}
