import random

import numpy as np
from sklearn.base import BaseEstimator
from sklearn.base import RegressorMixin

from . import binding


class XGPRegressor(BaseEstimator, RegressorMixin):

    def __init__(self, const_max=5, const_min=-5, funcs_string='sum,sub,mul,div', generations=30,
                 loss_metric='mae', max_height=6, min_height=3, n_pops=1, parsimony_coeff=0,
                 p_constant=0.5, p_terminal=0.3, pop_size=30, random_state=None, rounds=1,
                 tuning_generations=10):

        self.const_max = const_max
        self.const_min = const_min
        self.funcs_string = funcs_string
        self.generations = generations
        self.loss_metric = loss_metric
        self.max_height = max_height
        self.min_height = min_height
        self.n_pops = n_pops
        self.parsimony_coeff = parsimony_coeff
        self.p_constant = p_constant
        self.p_terminal = p_terminal
        self.pop_size = pop_size
        self.random_state = random_state
        self.rounds = rounds
        self.tuning_generations = tuning_generations

    def fit(self, X, y=None, **fit_params):

        self.program_str_ = binding.fit(
            X=X,
            y=y,
            X_names=fit_params.get('feature_names', ['X{}'.format(i) for i in range(X.shape[1])]),
            const_min=self.const_min,
            const_max=self.const_max,
            eval_metric_name=fit_params.get('eval_metric', self.loss_metric),
            funcs_string=self.funcs_string,
            generations=self.generations,
            loss_metric_name=self.loss_metric,
            max_height=self.max_height,
            min_height=self.min_height,
            n_pops=self.n_pops,
            parsimony_coeff=self.parsimony_coeff,
            p_constant=self.p_constant,
            p_terminal=self.p_terminal,
            pop_size=self.pop_size,
            rounds=self.rounds,
            seed=self.random_state if self.random_state else random.randrange(2 ** 64),
            tuning_generations=self.tuning_generations,
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
