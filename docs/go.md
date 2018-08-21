# Go

## Installation

Once you have [installed Go](https://golang.org/dl/), you can install XGO like any other Go package.

```sh
go get github.com/MaxHalford/xgp
```

## Usage

Although the full API is available on [godoc](https://godoc.org/github.com/MaxHalford/xgp), only a subset of it is relevant if all you want to do is train a program on a dataset.

### Instantiation

The core struct for learning in XGP is the `GP`. A `GP` encapsulates all the logic for generating, evaluating, and evolving programs. Although you can instantiate an `GP` directly, you can (and should) do it by instantiating a `GPConfig` and calling it's `NewGP` method. You can also use the `NewDefaultGPConfig` method to instantiate a `GPConfig` with the default values outlines in the [training parameters section](training-parameters.md). Even if you don't want to use the default values, it's a good idea to use `NewDefaultGPConfig` and then to set the fields you want to modify afterwards.

```go
var config = NewDefaultGPConfig()

config.LossMetric = metrics.Accuracy{}
config.Individuals = 42
config.Funcs = "cos,sin,exp"

var estimator = config.NewGP()
```

The `GPConfig` struct fields exactly match the ones indicated in the [training parameters section](training-parameters.md).

### Training

Once you have an `GP`, you're ready to call to it's `Fit` method to train it on a dataset. Here is the signature of the `Fit` method:

```go
func (est *GP) Fit(
    // Required arguments
    X [][]float64,
    Y []float64,
    // Optional arguments (can safely be nil)
    W []float64,
    XVal [][]float64,
    YVal []float64,
    WVal []float64,
    verbose bool,
) error
```

Just like for the CLI, the only required arguments to the `GP`'s `Fit` method are a matrix of features `X` and a list of targets `Y`. `W` can be used to weight the samples in `X` during program evaluation, which is particularly useful for higher-level learning algorithms such as [boosting](https://www.wikiwand.com/en/Boosting_(machine_learning)). One important thing to notice is that **`X` and `XVal` should be ordered column-wise**; that is `X[0]` should access the first column in the dataset, not the first row.

!!! warning
    For the while XGP does not handle categorical data. You should preemptively encode the categorical features in your dataset before feeding it to XGP. The recommended way is to use [label encoding](http://scikit-learn.org/stable/modules/preprocessing_targets.html#label-encoding) for ordinal data and [one-hot encoding](http://scikit-learn.org/stable/modules/preprocessing.html#encoding-categorical-features) for non-ordinal data.

!!! warning
    For the while XGP does not handle missing values.

Like the `val_set` argument in the CLI, `XVal`, `YVal`, and `WVal` can be used to track the performance of the best program on out-of-bag data. `notifyEvery` can be used to indicate at what frequency (in terms of genetic algorithm generations) progress should be displayed.

You can extract the best obtained `Program` with the `BestProgram` method.

```go
var best = gp.BestProgram()
```

Finally the `Fit` method returns an error which you should handle.

### Prediction

Once the `Fit` method has been called, the `Predict` method can be used to make predictions given a set of features.

```go
func (est GP) Predict(X [][]float64, predictProba bool) ([]float64, error)
```

The columns in `X` should be ordered in the same way as in the training set. The `proba` argument can be used to indicate if probabilities should be returned in the case of classification.

