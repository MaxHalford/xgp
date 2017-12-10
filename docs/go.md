# Go

## Installation

Once you have [installed Go](https://golang.org/dl/), you can install koza like any other Go package.

```sh
go get github.com/MaxHalford/koza
```

## Usage

Although the full API is available on [godoc](https://godoc.org/github.com/MaxHalford/koza), only a subset of it is relevant if all you want to do is train a program on a dataset.

### Instantiation

The core struct for learning in koza is the `Estimator`. An `Estimator` encapsulates all the logic for generating, evaluating, and evolving programs. Although you can instantiate an `Estimator` directly, you can (and should) do it by instantiating a `Config` and calling it's `NewEstimator` method. You can also use the `NewConfigWithDefaults` method to instantiate a `Config` with the default values outlines in the [training parameters section](training-parameters.md). Even if you don't want to use the default values, it's a good idea to use `NewConfigWithDefaults` and then to set the fields you want to modify afterwards.

```go
var config = NewConfigWithDefaults()

config.Individuals = 42
config.Funcs = "cos,sin,exp"

var estimator = config.NewEstimator()
```

The `Config` struct fields exactly match the ones indicated in the [training parameters section](training-parameters.md).

### Training

Once you have an `Estimator`, you're ready to call to it's `Fit` method to train it on a dataset. Here is the signature of the `Fit` method:

```go
func (est *Estimator) Fit(
    // Required arguments
    XTrain [][]float64,
    YTrain []float64,
    // Optional arguments
    WTrain []float64,
    XVal [][]float64,
    YVal []float64,
    WVal []float64,
    notifyEvery uint,
) error
```

Just like the in CLI, the only required arguments to the `Estimator`'s `Fit` are a matrix of features `XTrain` and a list of targets `YTrain`. `WTrain` can be used to weight the samples in `XTrain` during program evaluation, this is particularly useful for higher-level learning algorithms such as [boosting](https://www.wikiwand.com/en/Boosting_(machine_learning)). One important thing to not is that `XTrain` and `XVal` should be **ordered column-wise**; that is `X[0]` should access the first column in the dataset, not the first row.

!!! warning
    For the while koza does not handle missing values.

Like the `val_set` argument in the CLI, `XVal`, `YVal`, and `WVal` can be used to track the performance of the best program on out-of-bag data. `notifyEvery` can be used to indicate at what frequency (in terms of genetic algorithm generations) progress should be displayed.

Finally the `Fit` method returns an error which you should handle.

### Prediction

Once the `Fit` method has been called, then the `Predict` method can be called to make predictions on a dataset. Just like in the CLI, the columns in the test set should be ordered in the same way as in the training set. Here is the `Predict` method's signature:

```go
func (est Estimator) Predict(X [][]float64, predictProba bool) ([]float64, error)
```

Apart from the input dataset, the `predictProba` argument can be used to indicate if probabilities should be returned in the case of a classification task.

