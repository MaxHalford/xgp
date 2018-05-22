# Python

## Installation

Since version 1.5, Go code can be imported from Python as a dynamic-link library (DLL). This is what is done in the [XGP Python package](https://github.com/MaxHalford/xgp-python). As of now we aren't building Python wheels so you have to have Go and a GCC compiler installed on your machine. If this is the case you can install the package from PyPI:

```sh
>>> pip install xgp
```

This uses the [setuptools-golang project](https://github.com/asottile/setuptools-golang) to pull the needed Go dependencies and compile the DLL. In the near future the default behaviour be to use a pre-compiled DLL so you don't have to have Go and GCC installed.

## Usage

The XGP Python package exposes a scikit-learn API so you can use it in the same way as you would any other scikit-learn compliant code. Check out the [Python examples](https://github.com/MaxHalford/xgp-python/tree/master/examples) to get a general feel.

There are two estimators that you can use depending on if you're doing classification or regression, namely `XGPClassifier` and `XGPRegressor`. Both inherit from `XGPModel` and their only difference is that they have a different default loss metric. What's more `XGPClassifier` has a `predict_proba` method. After training you'll be able to access the `program_str_` to see the best program the estimator found.

The following snippet shows a basic usage example of `XGPClassifier`.

```python
from sklearn import datasets
from sklearn import metrics
from sklearn import model_selection
import xgp


X, y = datasets.load_breast_cancer(return_X_y=True)
X_train, X_test, y_train, y_test = model_selection.train_test_split(
    X,
    y,
    random_state=42
)

model = xgp.XGPClassifier(
    loss_metric='logloss',
    funcs='sum,sub,mul,div',
    n_individuals=500,
    n_generations=100,
    random_state=42,
    parsimony_coefficient=0.01
)

model.fit(X_train, y_train, eval_set=(X_test, y_test), verbose=True)

metric = metrics.log_loss
print('Train log-loss: {:.5f}'.format(metric(y_train, model.predict_proba(X_train))))
print('Test log-loss: {:.5f}'.format(metric(y_test, model.predict_proba(X_test))))
print('Best program:', model.program_str_)
```

This gives the following output:

```
Train log-loss: 0.217573
Test log-loss: 0.191963
Best program: sum(mul(X[0], mul(-4.774751043817239, X[7])), 3.8762056339039415)
```

The full list of parameters is available in the [training parameters section](training-parameters.md).

!!! warning
    For the while XGP does not handle categorical data. You should preemptively encode the categorical features in your dataset before feeding it to XGP. The recommended way is to use [label encoding](http://scikit-learn.org/stable/modules/preprocessing_targets.html#label-encoding) for ordinal data and [one-hot encoding](http://scikit-learn.org/stable/modules/preprocessing.html#encoding-categorical-features) for non-ordinal data.

!!! warning
    For the while xgp does not handle missing values.
