<div align="center">
  <!-- Logo -->
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vSLdt85rEf3SQUBkpuWfXOclyUY7rdZ7RBoTuNIyCc3-liSpurbL3i7QfrzWBFr2LfwTfoAf_1i4Qwe/pub?w=378&h=223"/>
</div>

<br/>

<div align="center">
  <!-- Documentation -->
  <a href="https://maxhalford.github.io/xgp">
    <img src="https://img.shields.io/website-up-down-green-red/http/shields.io.svg?label=documentation" alt="documentation" />
  </a>
  <!-- godoc -->
  <a href="https://godoc.org/github.com/MaxHalford/xgp">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="godoc" />
  </a>
  <!-- Build status -->
  <a href="https://travis-ci.org/MaxHalford/xgp">
    <img src="https://img.shields.io/travis/MaxHalford/gago/master.svg?style=flat-square" alt="build_status" />
  </a>
  <!-- Coverage status -->
  <a href="https://coveralls.io/github/MaxHalford/xgp?branch=master">
    <img src="https://coveralls.io/repos/github/MaxHalford/xgp/badge.svg?branch=master&style=flat-square" alt="coverage_status" />
  </a>
  <!-- License -->
  <a href="https://opensource.org/licenses/MIT">
    <img src="http://img.shields.io/:license-mit-ff69b4.svg?style=flat-square" alt="license"/>
  </a>
</div>

<br/>

XGP is a machine learning library for performing [symbolic regression](https://www.wikiwand.com/en/Symbolic_regression). It can be used both for regression and classification tasks. Please refer to [the documentation](https://maxhalford.github.io/xgp) for an in-depth introduction to symbolic regression.

## Interfaces

The core library is written in Go but it can be used in different ways:

- [Command-line interface (CLI)](https://maxhalford.github.io/xgp/cli/)
- [Go API](https://maxhalford.github.io/xgp/go/)
- [Python API](https://maxhalford.github.io/xgp/python/)

## Usage examples

### Command-line interface (CLI)

```sh
>>> xgp fit train.csv
>>> xgp predict test.csv
```

### Go

```go
package main

import "github.com/MaxHalford/xgp"

func main() {
    config := xgp.NewConfigWithDefaults()
    estimator := config.NewEstimator()

    estimator.Fit(XTrain, YTrain)
    yPred := estimator.Predict()
}
```

### Python

```python
import xgp

model = xgp.SymbolicRegressor()

model.fit(X_train, y_train)
y_pred = model.predict(X_test)
```

## Dependencies

The core xgp library has the following dependencies.

- [gonum/gonum](https://github.com/gonum/gonum) for [SIMD operations](https://www.wikiwand.com/en/SIMD) and calculating metrics
- [gosuri/uiprogress](https://github.com/gosuri/uiprogress) for displaying progress bars
- [MaxHalford/gago](https://github.com/MaxHalford/gago) for the genetic algorithm part

## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE) for more information.
