# Python usage

## Installation

Since version 1.5, Go code can be imported from Python as a dynamic-link library (DLL). This is what is done in the [XGP Python package](https://github.com/MaxHalford/xgp-python).

### Using a wheel

If you're using one of the following setups then you are in luck because a wheel is available. In other words you don't need to have Go and GCC installed.

|            | manylinux x86_64 |
|------------|:----------------:|
| Python 3.5 | ✅ |
| Python 3.6 | ✅ |

You need to have the `wheel` package installed.

```sh
>>> pip install wheel
```

Then you can install the wheel from PyPI.

```sh
>>> pip install xgp
```


### Compile it yourself

To compile the DLL you will need to have Go and GCC installed. Once this is done simply run:

```sh
>>> pip install --no-binary :all: xgp
```

This uses [setuptools-golang](https://github.com/asottile/setuptools-golang) to pull the needed Go dependencies and compile the DLL.

You can also build the DLL yourself with the following command. This is mostly for development purposes.

```sh
>>> go build -buildmode=c-shared -o xgp.so xgp/xgp.go
```

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
    flavor='vanilla',
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
```

This gives the following output:

```
Train log-loss: 0.217573
Test log-loss: 0.191963
```

The full list of parameters is available in the [training parameters section](training-parameters.md).

!!! warning
    For the while XGP does not handle categorical data. You should preemptively encode the categorical features in your dataset before feeding it to XGP. The recommended way is to use [label encoding](http://scikit-learn.org/stable/modules/preprocessing_targets.html#label-encoding) for ordinal data and [one-hot encoding](http://scikit-learn.org/stable/modules/preprocessing.html#encoding-categorical-features) for non-ordinal data.

!!! warning
    For the while XGP does not handle missing values.
