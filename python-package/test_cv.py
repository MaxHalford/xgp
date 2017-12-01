import time

from gplearn.genetic import SymbolicRegressor
import numpy as np
from sklearn import datasets
from sklearn import ensemble
from sklearn import linear_model
from sklearn import model_selection
from sklearn import pipeline
from sklearn import preprocessing
from sklearn import tree

from koza import regressor


X, Y = datasets.load_boston(return_X_y=True)

CV = model_selection.KFold(n_splits=5, random_state=42)

MODELS = {
    'Koza': regressor.SymbolicRegressor(
        population_size=300,
        n_generations=30,
        p_hoist_mutation=0.2,
        p_point_mutation=0.2,
        p_crossover=0.3,
        p_subtree_mutation=0.2,
        parsimony_coeff=0.001,
        random_state=42
    ),
    'gplearn': SymbolicRegressor(
        generations=30,
        p_crossover=0.3,
        p_hoist_mutation=0.2,
        p_point_mutation=0.2,
        p_subtree_mutation=0.2,
        population_size=300,
        random_state=2
    ),
    'Lasso': linear_model.Lasso(),
    'Ridge': linear_model.Ridge(),
    'Tree': tree.DecisionTreeRegressor(random_state=42),
    'Random forest': ensemble.RandomForestRegressor(random_state=42),
}


if __name__ == '__main__':

    left_space = max(len(name) for name in MODELS.keys()) + 1

    print(f'{"Model".rjust(left_space)} | Time (s.) |        MAE', )
    print(f'{"-" * left_space} + --------- + ------------------')
    for name, model in MODELS.items():
        t0 = time.time()
        scores = model_selection.cross_val_score(model, X=X, y=Y, scoring='neg_mean_absolute_error', cv=CV)
        mean_score = '{0:.4f}'.format(-np.mean(scores))
        std_score = '{0:.4f}'.format(np.std(scores))
        duration = '{0:.4f}'.format(time.time() - t0).rjust(9)
        print(f"{name.rjust(left_space)} | {duration} | {mean_score} (± {std_score})")


