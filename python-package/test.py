import numpy as np
from sklearn import datasets
from sklearn import linear_model
from sklearn import model_selection
from sklearn import pipeline
from sklearn import preprocessing

import regressor


X, y = datasets.load_boston(return_X_y=True)

models = {
    'XGP': pipeline.Pipeline([
        ('scale', preprocessing.StandardScaler()),
        ('xgp', regressor.XGPRegressor())
    ]),
    'Lasso': linear_model.Lasso(),
    'Ridge': linear_model.Ridge()
}

for name, model in models.items():
    scores = model_selection.cross_val_score(model, X=X, y=y, scoring='neg_mean_absolute_error')
    print(f'{name}: {-np.mean(scores)} (± {np.std(scores)})')
