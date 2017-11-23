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


X, y = datasets.load_boston(return_X_y=True)

cv = model_selection.KFold(n_splits=5, random_state=42)

models = {
    'Koza': regressor.SymbolicRegressor(
        population_size=300,
        n_generations=30,
        p_hoist_mutation=0.2,
        p_point_mutation=0.2,
        p_subtree_crossover=0.3,
        p_subtree_mutation=0.2,
        random_state=42
    ),
    'gplearn': SymbolicRegressor(
        generations=30,
        p_crossover=0.3,
        p_hoist_mutation=0.2,
        p_point_mutation=0.2,
        p_subtree_mutation=0.2,
        population_size=300,
        random_state=2,
        #verbose=1,
    ),
    'Linear': linear_model.LinearRegression(),
    'Lasso': linear_model.Lasso(),
    'Ridge': linear_model.Ridge(),
    'Tree': tree.DecisionTreeRegressor(random_state=42),
    'Random forest': ensemble.RandomForestRegressor(random_state=42),
}

for name, model in models.items():
    t0 = time.time()
    scores = model_selection.cross_val_score(model, X=X, y=y, scoring='neg_mean_absolute_error', cv=cv)
    print(f'{name}: {-np.mean(scores)} (± {np.std(scores)}) in {time.time() - t0} seconds')
