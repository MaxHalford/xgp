# Command-line interface (CLI)

## Examples

To generate the datasets for the following examples you need to have Python alongside with the [pandas](https://pandas.pydata.org/) and [scikit-learn](http://scikit-learn.org/stable/) libraries. The versions do not matter much.

For development you can replace `koza` with `go run main.go`. For example instead of running `koza fit examples/boston/train.csv --loss mae` you can run `go run main.go fit examples/boston/train.csv --loss mae`.

### Boston house prices

```sh
python examples/boston/create_datasets.py
koza fit examples/boston/train.csv --loss mae
koza score examples/boston/test.csv --eval mae
```

### gplearn polynomial example

This example is derived from [gplearn's documentation](http://gplearn.readthedocs.io/en/stable/examples.html). The function to approximate is a trivial and is only here for show.

```sh
python examples/gplearn/create_datasets.py
koza fit examples/gplearn/train.csv --loss mae
koza score examples/gplearn/test.csv --eval mae
```

### Titanic survivors

The data munging is adapted from this [Kaggle kernel](https://www.kaggle.com/scirpus/genetic-programming-lb-0-88).

```sh
python examples/titanic/create_datasets.py
koza fit examples/titanic/train.csv --loss mae
koza score examples/titanic/val.csv --eval mae
koza predict examples/titanic/test.csv examples/titanic/test_prediction.csv
```
