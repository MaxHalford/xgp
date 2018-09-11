<!-- This isn't a typo -->
#

<div align="center">
  <!-- Logo -->
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vSLdt85rEf3SQUBkpuWfXOclyUY7rdZ7RBoTuNIyCc3-liSpurbL3i7QfrzWBFr2LfwTfoAf_1i4Qwe/pub?w=378&h=223"/>
</div>

<div align="center">
  <!-- godoc -->
  <a href="https://godoc.org/github.com/MaxHalford/xgp">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="godoc" />
  </a>
  <!-- Build status -->
  <a href="https://travis-ci.org/MaxHalford/xgp">
    <img src="https://img.shields.io/travis/MaxHalford/eaopt/master.svg?style=flat-square" alt="build_status" />
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

XGP is a tool for performing [symbolic regression](https://www.wikiwand.com/en/Symbolic_regression) oriented towards machine learning. It can be used for classification and regression tasks. Please refer to the ["How it works" section](how-it-works.md) for an in-depth introduction to symbolic regression.

Symbolic regression is a heavy algorithm to run, hence a good implementation is paramount to make it usable in practice. XGP is written in Go, a compiled language which is both [fast](https://julialang.org/benchmarks/) and easy to write readable code with. The downside is that Go is not very practical for common data wrangling tasks. XGP is thus available as a CLI tool and can be imported by other languages (such as Python) to facililate it's usage and integration into data science pipelines.
