from sklearn import datasets

import regressor


X, y = datasets.load_boston(return_X_y=True)

clf = regressor.XGPRegressor()
clf.fit(X, y)
y_pred = clf.predict(X)

print(y_pred)
