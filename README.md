<div align="center">
  <!-- Logo -->
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vSLdt85rEf3SQUBkpuWfXOclyUY7rdZ7RBoTuNIyCc3-liSpurbL3i7QfrzWBFr2LfwTfoAf_1i4Qwe/pub?w=378&h=223"/>
</div>

<br/>

<div align="center">
  <!-- Documentation -->
  <a href="https://maxhalford.github.io/koza">
    <img src="https://img.shields.io/website-up-down-green-red/http/shields.io.svg?label=documentation" alt="documentation" />
  </a>
  <!-- godoc -->
  <a href="https://godoc.org/github.com/MaxHalford/koza">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="godoc" />
  </a>
  <!-- Build status -->
  <a href="https://travis-ci.org/MaxHalford/koza">
    <img src="https://img.shields.io/travis/MaxHalford/gago/master.svg?style=flat-square" alt="build_status" />
  </a>
  <!-- License -->
  <a href="https://opensource.org/licenses/MIT">
    <img src="http://img.shields.io/:license-mit-ff69b4.svg?style=flat-square" alt="license"/>
  </a>
</div>

<br/>

<div align="center">A <a href="https://www.wikiwand.com/en/Symbolic_regression">symbolic regression</a> tool written in Go with bindings to other languages</div>

<br/>

:warning: koza is still in active development phase.

koza is a tool for performing symbolic regression oriented towards machine learning. It can be used for regression and classification tasks. Please refer to [the documentation](https://maxhalford.github.io/koza) for an in-depth introduction to symbolic regression.

<br/>
<div align="center">
  <a href="https://asciinema.org/a/x6t8d5PZ4Td6iDoAa4IXeK7IB">
    <img src="https://asciinema.org/a/x6t8d5PZ4Td6iDoAa4IXeK7IB.png" width="60%" />
  </a>
</div>
<br/>

## Interfaces

The core library is written in Go but can be used in different ways.

- [Command-line interface (CLI)](https://maxhalford.github.io/koza/cli/)
- [Go API](https://maxhalford.github.io/koza/go/)
- [Python API](https://maxhalford.github.io/koza/python/)

## Usage examples

### Command-line interface (CLI)

```sh
>>> koza fit train.csv
>>> koza predict test.csv
```

### Go

```go
package main

import "github.com/MaxHalford/koza"

func main() {
    config := koza.NewConfigWithDefaults()
    estimator := config.NewEstimator()

    estimator.Fit(XTrain, YTrain)
    yPred := estimator.Predict()
}
```

### Python

```python
import koza

model = koza.SymbolicRegressor()

model.fit(X_train, y_train)
y_pred = model.predict(X_test)
```

## Thanks

koza uses the following projects which are a joy to work with.

- [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) for displaying parameters in a pretty way
- [gonum/gonum](https://github.com/gonum/gonum) for [SIMD operations](https://www.wikiwand.com/en/SIMD) and calculating metrics
- [spf13/cobra](https://github.com/spf13/cobra) for building the CLI
- [kniren/gota](https://github.com/kniren/gota) for manipulating dataframes from the CLI
- [gosuri/uiprogress](https://github.com/gosuri/uiprogress) for displaying progress bars in the CLI
- [mkdocs/mkdocs](https://github.com/mkdocs/mkdocs/) and [squidfunk/mkdocs-material](https://github.com/squidfunk/mkdocs-material) for building the documentation

## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE) for more information.
