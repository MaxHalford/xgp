from sklearn.base import BaseEstimator
from sklearn.base import ClassifierMixin

from . import binding


class XGPClassifier(BaseEstimator, ClassifierMixin):

    def __init__(self, generations=10, tuning_generations=10):
        self.generations = generations
        self.tuning_generations = tuning_generations

    def fit(self, X, y=None, **fit_params):
        return binding.fit(
            X=X,
            y=y,
            metric_name='mse',
            generations=self.generations,
            tuning_generations=self.tuning_generations
        )

    def predict(self, X):
        return binding.predict(X)
