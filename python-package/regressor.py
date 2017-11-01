from sklearn.base import BaseEstimator
from sklearn.base import RegressorMixin

import binding


class XGPRegressor(BaseEstimator, RegressorMixin):

    def __init__(self, const_max=5, const_min=-5, eval_metric_name='mae', feature_names=None,
                 funcs_string='sum,sub,mul,div', generations=30, loss_metric_name='mae',
                 max_height=6, min_height=3, parsimony_coeff=0, p_constant=0.5, p_terminal=0.3,
                 rounds=1, tuning_generations=10, verbose=True):

        self.const_max = const_max
        self.const_min = const_min
        self.eval_metric_name = eval_metric_name
        self.feature_names = feature_names
        self.funcs_string = funcs_string
        self.generations = generations
        self.loss_metric_name = loss_metric_name
        self.max_height = max_height
        self.min_height = min_height
        self.parsimony_coeff = parsimony_coeff
        self.p_constant = p_constant
        self.p_terminal = p_terminal
        self.rounds = rounds
        self.tuning_generations = tuning_generations
        self.verbose = verbose

    def fit(self, X, y=None, **fit_params):

        if self.feature_names is None:
            self.feature_names = ['X{}'.format(i) for i in range(X.shape[1])]

        return binding.fit(
            X=X,
            y=y,
            X_names=self.feature_names,
            const_min=self.const_min,
            const_max=self.const_max,
            eval_metric_name=self.eval_metric_name,
            funcs_string=self.funcs_string,
            generations=self.generations,
            loss_metric_name=self.loss_metric_name,
            max_height=self.max_height,
            min_height=self.min_height,
            parsimony_coeff=self.parsimony_coeff,
            p_constant=self.p_constant,
            p_terminal=self.p_terminal,
            rounds=self.rounds,
            tuning_generations=self.tuning_generations,
            verbose=self.verbose
        )

    def predict(self, X):
        return binding.predict(X=X, predict_proba=False)
