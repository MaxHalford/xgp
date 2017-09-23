<div align="center">
  <!-- Logo -->
  <img src="https://docs.google.com/drawings/d/1en_XKo3L65RCiFtu2ftutXYpPE3DO7SBW3qLL36Rdg4/pub?w=389&h=227"/>
</div>

xgp is a machine learning tool based on [genetic programming](https://www.wikiwand.com/en/Genetic_programming) which can be used for both classification and regression problems.

## Installation

- [Command-line interface (CLI)](cmd/README.md)
- [Python package](python-package/README.md)

You can also simply install the `xgp` Go package and use it within your Go code.

```sh
go get -u github.com/MaxHalford/xgp
```

The following Go dependencies will be installed alongside:

- [MaxHalford/gago](https://github.com/MaxHalford/gago)

## Usage

If you want to use xgp in a data science pipeline then you probably want to use it in one of the following ways:

- [Command-line interface (CLI)](cmd/README.md)
- [Python package](python-package/README.md)

However, you can also use xgp inside your Go code as you would do with any other library. Although the full API is available on godoc, you will (and should) be using the following methods most of the time; which is in fact what is done by the CLI and the other languages packages.

```go
err := Estimator.Fit(X [][]float64, Y []float64)
```

```go
yPred, err := Estimator.Predict(X [][]float64)
```
