from sklearn import datasets

from xgp import classifier


X, y = datasets.load_diabetes(return_X_y=True)

clf = classifier.XGPClassifier()
clf.fit(X, y)
yPred = clf.predict(X)

print(yPred)
