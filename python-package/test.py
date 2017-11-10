import numpy as np
from sklearn import datasets
from sklearn import ensemble
from sklearn import linear_model
from sklearn import model_selection
from sklearn import pipeline
from sklearn import preprocessing
from sklearn import tree

import regressor


X, y = datasets.load_boston(return_X_y=True)

cv = model_selection.KFold(n_splits=3, random_state=42)

models = {
    'Random forest': ensemble.RandomForestRegressor(random_state=42),
    'XGP': regressor.XGPRegressor(random_state=42),
    'Lasso': linear_model.Lasso(),
    'Ridge': linear_model.Ridge(),
    'Tree': tree.DecisionTreeRegressor(random_state=42)
}

for name, model in models.items():
    scores = model_selection.cross_val_score(model, X=X, y=y, scoring='neg_mean_absolute_error', cv=cv)
    print(f'{name}: {-np.mean(scores)} (± {np.std(scores)})')
