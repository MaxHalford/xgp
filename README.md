<div align="center">
  <!-- Logo -->
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vSLdt85rEf3SQUBkpuWfXOclyUY7rdZ7RBoTuNIyCc3-liSpurbL3i7QfrzWBFr2LfwTfoAf_1i4Qwe/pub?w=378&h=223"/>
</div>

xgp is a machine learning tool for classification and regression based on [genetic programming](https://www.wikiwand.com/en/Genetic_programming).

## Installation

If you want to use xgp in a data science pipeline then you probably want to use it in one of the following ways:

- [Command-line interface (CLI)](cmd/README.md)
- [Python package](python-package/README.md)

You can also simply install the `xgp` Go package and use it within your Go code.

```sh
go get -u github.com/MaxHalford/xgp
```

The following Go dependencies will be installed alongside:

- [gonum/gonum](https://github.com/gonum/gonum)
- [MaxHalford/gago](https://github.com/MaxHalford/gago)

## Go usage

Although the full API is available on godoc, you will (and should) be using the `Fit` and `Predict` methods from `Estimator` struct; which is in fact what is done by the CLI and client libraries.

```go
var estimator = Estimator{} // Set parameters here
err := estimator.Fit(X [][]float64, Y []float64)
yPred, err := estimator.Predict(X [][]float64)
```

The `Estimator` struct fields has many fields you can set, most of which are accessible via the CLI and client libraries.
