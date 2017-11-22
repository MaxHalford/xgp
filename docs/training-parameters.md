# Training parameters

## Quick overview

| CLI             | Go              | Python          | Description                                  | Default |
|-----------------|-----------------|-----------------|----------------------------------------------|---------|
| `const_max`     | `constMax`      | `const_max`     | Maximum value used for generating constants. | -5 |
| `const_min`     | `constMin`      | `const_min`     | Minimum value used for generating constants. |  5 |
| `eval` | `evalMetric` | `eval_metric` | Evaluation metric used for monitoring progress. Has to be passed in the `fit` method for Python. | Same as `loss_metric` |
| `funcs`   | `funcs`   | `funcs_string`   | Functions that can be used to generate programs.         | `sum,sub,mul,div` |
| `loss`   | `lossMetric`   | `loss_metric`   | Loss metric used for evaluating programs. This determines if the task is a classification or a regression one. | `mae` |
| `max_height` | `maxHeight` | `max_height` | Maximum program height used in ramped half-and-half initialization. | 6 |
| `min_height` | `minHeight` | `min_height` | Minimum program height used in ramped half-and-half initialization. | 3 |
| `n_generations`   | `generations`   | `n_generations`   | Number of generations used by the genetic algorithm.         | 30 |
| `n_pops`   | `nPops`   | `n_populations`   | Number of populations used by the genetic algorithm.         | 1 |
| `parsimony`   | `parsimonyCoeff`   | `parsimony_coeff`   | Parsimony coefficient used for evaluating programs.         | 0 |
| `p_constant`   | `pConstant`   | `p_constant`   | Probability of generating a constant during ramped half-and-half initialization. | 0.5 |
| `p_full`   | `pFull`   | `p_full`   | Probability of using full initialization during ramped half-and-half initialization. | 0.5 |
| `p_hoist_mut`   | `pHoistMutation`   | `p_hoist_mutation`   | Probability of applying hoist mutation.         | 0.2 |
| `p_point_mut`   | `pPointMutation`   | `p_point_mutation`   | Probability of applying point mutation.         | 0.2 |
| `p_point_cross`   | `pSubTreeCrossover`   | `p_subtree_crossover`   | Probability of applying sub-tree crossover.         | 0.3 |
| `p_point_lut`   | `pSubTreeMutation`   | `p_subtree_mutation`   | Probability of applying sub-tree mutation.         | 0.2 |
| `p_terminal`   | `pTerminal`   | `p_terminal`   | Probability of generating a terminal node during ramped half-and-half initialization. | 0.5 |
| `pop_size`   | `pPopSize`   | `population_size`   | Probability of generating a terminal node during ramped half-and-half initialization. | 0.5 |
| `rounds`   | `rounds`   | `n_rounds`   | Number of rounds used for boosting (not yet implemented). | 1 |
| `seed`   | `seed`   | `random_state`   | Seed used for generating random numbers. | Random |
| `verbose`   | `verbose`   | `verbose`   | Indicates if progress should be monitored or not. | False |


## Detailed walkthrough
