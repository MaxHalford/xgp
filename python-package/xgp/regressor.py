import random

import numpy as np
from sklearn.base import BaseEstimator
from sklearn.base import RegressorMixin

from . import binding


class XGPRegressor(BaseEstimator, RegressorMixin):

    def __init__(self, const_max=5, const_min=-5, funcs_string='sum,sub,mul,div', loss_metric='mae',
                 max_height=6, min_height=3, n_generations=30, n_populations=1, parsimony_coeff=0,
                 p_constant=0.5, p_hoist_mutation=0.2, p_point_mutation=0.2,
                 p_subtree_crossover=0.3, p_subtree_mutation=0.2, p_terminal=0.5,
                 population_size=30, random_state=None, n_rounds=1, tuning_n_generations=0):

        self.const_max = const_max
        self.const_min = const_min
        self.funcs_string = funcs_string
        self.loss_metric = loss_metric
        self.max_height = max_height
        self.min_height = min_height
        self.n_generations = n_generations
        self.n_populations = n_populations
        self.parsimony_coeff = parsimony_coeff
        self.p_constant = p_constant
        self.p_hoist_mutation = p_hoist_mutation
        self.p_point_mutation = p_point_mutation
        self.p_subtree_crossover = p_subtree_crossover
        self.p_subtree_mutation = p_subtree_mutation
        self.p_terminal = p_terminal
        self.population_size = population_size
        self.random_state = random_state
        self.n_rounds = n_rounds
        self.tuning_n_generations = tuning_n_generations

    def fit(self, X, y=None, **fit_params):

        self.program_str_ = binding.fit(
            X=X,
            y=y,
            X_names=fit_params.get('feature_names', ['X{}'.format(i) for i in range(X.shape[1])]),
            const_min=self.const_min,
            const_max=self.const_max,
            eval_metric_name=fit_params.get('eval_metric', self.loss_metric),
            funcs_string=self.funcs_string,
            loss_metric_name=self.loss_metric,
            max_height=self.max_height,
            min_height=self.min_height,
            n_generations=self.n_generations,
            n_populations=self.n_populations,
            parsimony_coeff=self.parsimony_coeff,
            p_constant=self.p_constant,
            p_hoist_mutation=self.p_hoist_mutation,
            p_point_mutation=self.p_point_mutation,
            p_subtree_crossover=self.p_subtree_crossover,
            p_subtree_mutation=self.p_subtree_mutation,
            p_terminal=self.p_terminal,
            population_size=self.population_size,
            n_rounds=self.n_rounds,
            seed=self.random_state if self.random_state else random.randrange(2 ** 16),
            tuning_n_generations=self.tuning_n_generations,
            verbose=fit_params.get('verbose', False)
        )

        self.program_eval_ = lambda X: eval(self.program_str_)

        return self

    def predict(self, X):
        y_pred = self.program_eval_(X)

        # In case the program is a single constant it has to be converted to an array
        if isinstance(y_pred, float):
            y_pred = np.array([y_pred] * len(X))

        return y_pred
