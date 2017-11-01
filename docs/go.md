# Go

## Installation

Once you have [installed Go](https://golang.org/dl/), you can install the `xgp` package with the following one-liner.

```sh
go get -u github.com/MaxHalford/xgp
```

## Usage

Although the full API is available on godoc, you will (and should) be using the `Fit` and `Predict` methods from `Estimator` struct; which is in fact what is done by the CLI and the client libraries.

```go
estimator := Estimator{} // Set parameters here
err := estimator.Fit(X [][]float64, Y []float64)
yPred, err := estimator.Predict(X [][]float64)
```

The `Estimator` struct fields has many fields you can set, most of which are accessible via the CLI and client libraries. Check the [parameters page](parameters.md) for more information.
