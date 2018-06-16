# Training parameters

## Overview

The following tables gives an overview of all the parameters that can be used for training XGP. The defaults are the same regardless of where you're using XGP from (please [open an issue](https://github.com/MaxHalford/xgp/issues/new) if you notice any descrepancies). The values indicated for Go are the ones that can be passed to a `Config` struct. For Python some parameters have to be passed in the `fit` method.

### Learning parameters

| Name | CLI | Go | Python | Default value |
|------|-----|----|--------|---------------|
| Loss metric; is used to if the task is classification or regression | `loss` | `LossMetricName` | `loss_metric` | mae (for Python `XGPClassifier` defaults to logloss) |
| Evaluation metric | `eval` | `EvalMetricName` | `eval_metric` (in `fit`) | Same as loss metric |
| Parsimony coefficient | `parsimony` | `ParsimonyCoefficient` | `parsimony_coeff` | 0.00001 |
| Polish the best program | `polish` | `PolishBest` | `polish_best` | true |

Because XGP doesn't require the loss metric to be differentiable you can use any loss metric available. If you don't specify an evaluation metric then it will default to using the loss metric.

### Function parameters

| Name | CLI | Go | Python | Default value |
|------|-----|----|--------|---------------|
| Authorized functions | `funcs` | `Funcs` | `funcs` | sum,sub,mul,div |
| Constant minimum | `const_min` | `ConstMin` | `const_min` | -5 |
| Constant maximum | `const_max` | `ConstMax` | `const_max` | 5 |
| Constant probability  | `p_const` | `PConst` | `p_const` | 0.5 |
| Full initialization probability  | `p_full` | `PFull` | `p_full` | 0.5 |
| Terminal probability  | `p_leaf` | `PLeaf` | `p_leaf` | 0.3 |
| Minimum height | `min_height` | `MinHeight` | `min_height` | 3 |
| Maximum height | `max_height` | `MaxHeight` | `max_height` | 5 |

These parameters are used to generate the initial set of programs. They will also be used to generate new programs for subtree mutation. XGP uses ramped half-and-half initialization; the full initialization probability determines the probability of using full initialization and consequently the probability of using grow initialization.

### Genetic algorithm parameters

| Name | CLI | Go | Python | Default value |
|------|-----|----|--------|---------------|
| Number of populations | `pops` | `NPopulations` | `n_populations` | 1 |
| Number of individuals per population | `indis` | `NIndividuals` | `n_individuals` | 50 |
| Number of generations | `gens` | `NGenerations` | `n_generations` | 30 |
| Hoist mutation probability | `p_hoist_mut` | `PHoistMutation` | `p_hoist_mutation` | 0.1 |
| Subtree mutation probability | `p_sub_mut` | `PSubtreeMutation` | `p_sub_tree_mutation` | 0.1 |
| Point mutation probability | `p_point_mut` | `PPointMutation` | `p_point_mutation` | 0.1 |
| Point mutation rate | `point_mut_rate` | `PointMutationRate` | `point_mutation_rate` | 0.3 |
| Subtree crossover probability | `p_sub_cross` | `PSubtreeCrossover` | `p_sub_tree_crossover` | 0.5 |

### Other parameters

| Name | CLI | Go | Python | Default value |
|------|-----|----|--------|---------------|
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

Code-wise the operators are all located in the `op` subpackage, of which the goal is to provide fast implementations for each operator. For the while the only accelerations that exist are the ones for the sum and the division which use assembly implementations made available by [`gonum/floats`](https://godoc.org/gonum.org/v1/gonum/floats).

| Name | Arity | Short name | Go struct |
|------|-------|------------|-----------|
| Absolute value | 1 | abs | `Abs` |
| Addition | 2 | add | `Add` |
| Cosine | 1 | cos | `Cos` |
| Division | 2 | div | `Div` |
| Inverse | 1 | inv | `Inv` |
| Maximum | 2 | max | `Max` |
| Minimum | 2 | min | `Min` |
| Multiplication | 2 | mul | `Mul` |
| Negative value | 1 | neg | `Neg` |
| Sine | 1 | sin | `Sin` |
| Square | 2 | square | `Square` |
| Subtraction | 2 | sub | `Sub` |

Safe-division is used, meaning that if a denominator is 0 then the result will default to 1.
