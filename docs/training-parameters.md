# Training parameters

## Quick overview

The following table gives an overview of all the parameters that can be used for training koza. The defaults are the same regardless of the API. The values indicated for Go are the ones that can be passed to a `Config` struct.

| Name | CLI | Go | Python | Default |
|------|-----|----|--------|---------|
| Constant minimum | const_min | ConstMin | const_min | -5 |
| Constant maximum | const_max | ConstMax | const_max | 5 |
| Evaluation metric | eval | EvalMetricName | eval_metric (in `fit`) | mae |
| Loss metric | loss | LossMetricName | loss_metric | Same as loss metric |
| Authorized functions | funcs | Funcs | funcs | sum,sub,mul,div |
| Minimum height | min_height | MinHeight | min_height | 3 |
| Maximum height | max_height | MaxHeight | max_height | 5 |
| Number of populations | pops | NPopulations | n_populations | 1 |
| Number of individuals per population | indis | NIndividuals | n_individuals | 50 |
| Number of generations | gens | NGenerations | n_generations | 30 |
| Number of tuning generations | tune_gens | NTuningGenerations | n_tuning_generations | 0 |
| Constant probability  | p_const | PConstant | p_const | 0.5 |
| Full initialization probability  | p_full | PFull | p_full | 0.5 |
| Terminal probability  | p_terminal | PTerminal | p_terminal | 0.3 |
| Hoist mutation probability | p_hoist_mut | PHoistMutation | p_hoist_mutation | 0.1 |
| Sub-tree mutation probability | p_sub_mut | PSubTreeMutation | p_sub_tree_mutation | 0.1 |
| Point mutation probability | p_point_mut | PPointMutation | p_point_mutation | 0.1 |
| Point mutation rate | point_mut_rate | PointMutationRate | point_mutation_rate | 0.3 |
| Sub-tree crossover probability | p_sub_cross | PSubTreeCrossover | p_sub_tree_crossover | 0.5 |
| Parsimony coefficient | parsimony | ParsimonyCoefficient | parsimony_coeff | 0 |
| Random number seed | seed | Seed | seed | Random |

## Detailed walkthrough

TODO
