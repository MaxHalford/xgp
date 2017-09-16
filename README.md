<div align="center">
  <!-- Logo -->
  <img src="https://docs.google.com/drawings/d/1en_XKo3L65RCiFtu2ftutXYpPE3DO7SBW3qLL36Rdg4/pub?w=389&h=227"/>
</div>

xgp is a machine learning tool based on [genetic programming](https://www.wikiwand.com/en/Genetic_programming) which can be used for both classification and regression problems.

## To do

- Boosting command-line
- Create lookup tables for operators that take time to evaluate at runtime
- http://www.genetic-programming.com/gpquadraticexample.html
- http://www.genetic-programming.com/gpflowchart.html
- [Classification Strategies for Image Classification in Genetic Programming](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.475.3010&rep=rep1&type=pdf)
- [Multi-class overview](http://dynamics.org/~altenber/UH_ICS/EC_REFS/GP_REFS/IEEE/CEC2001/395.pdf)
- [Sampling data for fitness evaluation](http://eplex.cs.ucf.edu/papers/morse_gecco16.pdf)
- http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.475.3010&rep=rep1&type=pdf
- http://cswww.essex.ac.uk/staff/poli/gp-field-guide/42StepbyStepSampleRun.html
- Consider parsimony for generalization

## Prerequisites

First, [install Go](https://golang.org/dl/), set your `GOPATH`, and make sure `$GOPATH/bin` is on your `PATH`.

```sh
brew install go # If using homebrew

# Put these in your .bash_profile or .zshrc
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

## Installation

### Go library

```sh
go get -u github.com/MaxHalford/xgp
```

This will install the xgp package. The following Go dependencies will be installed alongside:

- [MaxHalford/gago](https://github.com/MaxHalford/gago)

### CLI

```sh
go get -u github.com/MaxHalford/xgp/cmd
```

This will install the xgp package together with the CLI. The following Go dependencies will be installed alongside:

- [fatih/color](https://github.com/fatih/color)
- [gosuri/uiprogress](https://github.com/gosuri/uiprogress)
- [MaxHalford/gago](https://github.com/MaxHalford/gago)
- [spf13/cobra](https://github.com/spf13/cobra)
- [spf13/viper](https://github.com/spf13/viper)

### Python package
