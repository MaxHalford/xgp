# Training parameters

## Quick overview

The following table gives an overview of all the parameters that can be used for training XGP. The defaults are the same regardless of the API. The values indicated for Go are the ones that can be passed to a `Config` struct.

| Name | CLI | Go | Python | Default |
|------|-----|----|--------|---------|
| Constant minimum | `const_min` | `ConstMin` | `const_min` | -5 |
| Constant maximum | `const_max` | `ConstMax` | `const_max` | 5 |
| Evaluation metric | `eval` | `EvalMetricName` | `eval_metric` (in `fit`) | mae |
| Loss metric | `loss` | `LossMetricName` | `loss_metric` | Same as loss metric |
| Authorized functions | `funcs` | `Funcs` | `funcs` | sum,sub,mul,div |
| Minimum height | `min_height` | `MinHeight` | `min_height` | 3 |
| Maximum height | `max_height` | `MaxHeight` | `max_height` | 5 |
| Number of populations | `pops` | `NPopulations` | `n_populations` | 1 |
| Number of individuals per population | `indis` | `NIndividuals` | `n_individuals` | 50 |
| Number of generations | `gens` | `NGenerations` | `n_generations` | 30 |
| Number of tuning generations | `tune_gens` | `NTuningGenerations` | `n_tuning_generations` | 0 |
| Constant probability  | `p_const` | `PConstant` | `p_const` | 0.5 |
| Full initialization probability  | `p_full` | `PFull` | `p_full` | 0.5 |
| Terminal probability  | `p_terminal` | `PTerminal` | `p_terminal` | 0.3 |
| Hoist mutation probability | `p_hoist_mut` | `PHoistMutation` | `p_hoist_mutation` | 0.1 |
| Sub-tree mutation probability | `p_sub_mut` | `PSubTreeMutation` | `p_sub_tree_mutation` | 0.1 |
| Point mutation probability | `p_point_mut` | `PPointMutation` | `p_point_mutation` | 0.1 |
| Point mutation rate | `point_mut_rate` | `PointMutationRate` | `point_mutation_rate` | 0.3 |
| Sub-tree crossover probability | `p_sub_cross` | `PSubTreeCrossover` | `p_sub_tree_crossover` | 0.5 |
| Parsimony coefficient | `parsimony` | `ParsimonyCoefficient` | `parsimony_coeff` | 0 |
| Random number seed | `seed` | `Seed` | `seed` | Random |

## Loss metrics

Genetic programming directly minimises a loss metric. Because the optimization is done with a genetic algorithm the loss metric doesn't have to be differentiable. Whether the task is classification or regression is thus determined from the loss metric. This is similar to how XGBoost and LightGBM handle things.

Each loss metric has a short name that you can use whether you are using the CLI, Go, or Python. You can also use these short names to evaluate the performance of the model. For example you might want to optimise the ROC AUC while also keeping track of the accuracy.

| Name | Short name | Task |
|------|------------|------|
| Logloss | logloss | Classification |
| Accuracy | accuracy | Classification |
| Precision | precision | Classification |
| Recall | recall | Classification |
| F1-score | f1 | Classification |
| ROC AUC | roc_auc | Classification |
| Mean absolute error | mae | Regression |
| Mean squared error | mse | Regression |
| Root mean squared error | rmse | Regression |
| R2 | r2 | Regression |
| Absolute Pearson correlation | pearson | Regression |

## Operators

The following table lists all the available operators. Regardless of from where it is being used from, functions should be passed to XGP by concatenating the short names of the functions with a comma. For example to use the natural logarithm and the multiplication use `log,mul`.

Code-wise the operators are all located in the `op` sub-package, of which the goal is to provide fast implementations for each operator. For the while the only accelerations that exist are the ones for the sum and the division which use assembly implementations made available by [gonum/floats](https://godoc.org/gonum.org/v1/gonum/floats).

| Name | Arity | Short name | Go struct | Assembly code |
|------|-------|------------|---------------|---------------|
| Cosine | 1 | cos | `Cos` | ✗ |
| Sine | 1 | sin | `Sin` | ✗ |
| Natural logarithm | 1 | log | `Log` | ✗ |
| Exponential | 1 | exp | `Exp` | ✗ |
| Maximum | 2 | max | `Max` | ✗ |
| Minimum | 2 | min | `Min` | ✗ |
| Sum | 2 | sum | `Sum` | ✔ |
| Subtraction | 2 | sub | `Sub` | ✗ |
| Division | 2 | div | `Div` | ✔ |
| Multiplication | 2 | mul | `Mul` | ✗ |
| Power | 2 | pow | `Pow` | ✗ |

