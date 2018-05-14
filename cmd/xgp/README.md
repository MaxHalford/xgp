# Command-line interface (CLI)

If you're looking for documentation then please refer to the [website](https://maxhalford.github.io/xgp/cli/).

## Examples

To generate the datasets for the following examples you need to have Python alongside with [pandas](https://pandas.pydata.org/) and [scikit-learn](http://scikit-learn.org/stable/). The versions do not really matter.

For development you can replace `xgp` with `go run main.go`. For example instead of running `xgp fit examples/boston/train.csv --loss mae` you can run `go run main.go fit examples/boston/train.csv --loss mae`.

### Boston house prices

```sh
>>> python examples/boston/create_datasets.py
>>> xgp fit examples/boston/train.csv --val examples/boston/test.csv --loss mae --seed 42 --indis 50 --gens 30
```

### gplearn polynomial example

This example is derived from [gplearn's documentation](http://gplearn.readthedocs.io/en/stable/examples.html). The function to approximate is a trivial and is only here for show.

```sh
>>> python examples/gplearn/create_datasets.py
>>> xgp fit examples/gplearn/train.csv --loss mae
>>> xgp score examples/gplearn/test.csv --eval mae
```

### Titanic survivors

The data munging is adapted from this [Kaggle kernel](https://www.kaggle.com/scirpus/genetic-programming-lb-0-88). Download the data from [here](https://www.kaggle.com/c/titanic/data) and put it in the `examples/titanic/kaggle` directory.

```sh
>>> python examples/titanic/create_datasets.py
>>> xgp fit examples/titanic/train.csv --loss logloss --eval accuracy --val examples/titanic/val.csv --target Survived --ignore PassengerId --parsimony 0.001 --gens 64 --indis 256 --funcs sum,sub,mul,div,cos,sin,min,max
>>> xgp predict examples/titanic/test.csv --output examples/titanic/submission.csv --keep PassengerId --target Survived
```
