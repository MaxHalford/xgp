import time

from gplearn.genetic import SymbolicRegressor
import numpy as np
import pandas as pd
from sklearn import ensemble
from sklearn import linear_model
from sklearn import metrics
from sklearn import pipeline
from sklearn import preprocessing
from sklearn import tree

from xgp import regressor


train = pd.read_csv('../examples/boston/train.csv')
test = pd.read_csv('../examples/boston/test.csv')

X_train = train.drop('y', axis='columns').values
y_train = train['y'].values
X_test = test.drop('y', axis='columns').values
y_test = test['y'].values

models = {
    'XGP': regressor.XGPRegressor(
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
        verbose=1,
    ),
    'Linear': linear_model.LinearRegression(),
    'Lasso': linear_model.Lasso(),
    'Ridge': linear_model.Ridge(),
    'Tree': tree.DecisionTreeRegressor(random_state=42),
    'Random forest': ensemble.RandomForestRegressor(random_state=42),
}

for name, model in models.items():
    t0 = time.time()
    model.fit(X_train, y_train)
    y_pred = model.predict(X_test)
    test_score = metrics.mean_absolute_error(y_test, y_pred)
    print(f'{name}: {test_score} in {time.time() - t0} seconds')
