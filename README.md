xgp is a machine learning tool based on [genetic programming](https://www.wikiwand.com/en/Genetic_programming) which can be used for classification and regression problems.

## To do

- Create lookup tables for operators that time to evaluate at runtime
- http://www.genetic-programming.com/gpquadraticexample.html
- http://www.genetic-programming.com/gpflowchart.html
- [Classification Strategies for Image Classification in Genetic Programming](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.475.3010&rep=rep1&type=pdf)
- [Multi-class overview](http://dynamics.org/~altenber/UH_ICS/EC_REFS/GP_REFS/IEEE/CEC2001/395.pdf)
- [Sampling data for fitness evaluation](http://eplex.cs.ucf.edu/papers/morse_gecco16.pdf)
- http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.475.3010&rep=rep1&type=pdf
- http://cswww.essex.ac.uk/staff/poli/gp-field-guide/42StepbyStepSampleRun.html
- "Simplify" an AST by reducing non-necessary branches (no variables)
- Consider parsimony for generalization

## Prerequisites

First, [install Go](https://golang.org/dl/), set your `GOPATH`, and make sure `$GOPATH/bin` is on your `PATH`.

```sh
brew install go # If using homebrew

# Put these in .bash_profile or .zshrc
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

## Installation

```
go get -u github.com/MaxHalford/xgp/...
```

This will the xgp package together with the CLI. The following Go dependencies will be installed:

- [MaxHalford/gago](https://github.com/MaxHalford/gago)
- [urfave/cli](https://github.com/urfave/cli)

## Development

```sh
go run cmd/xgp/*.go fit examples/regression/train.csv -tc y && go run cmd/xgp/*.go predict examples/regression/test.csv -tc y
```
