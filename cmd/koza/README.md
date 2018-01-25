# Command-line interface (CLI)

## Examples

To generate the datasets for the following examples you need to have Python alongside with the [pandas](https://pandas.pydata.org/) and [scikit-learn](http://scikit-learn.org/stable/) libraries. The versions do not matter much.

For development you can replace `koza` with `go run main.go`. For example instead of running `koza fit examples/boston/train.csv --loss mae` you can run `go run main.go fit examples/boston/train.csv --loss mae`.

```sh
go run main.go fit examples/restaurants/train.csv --loss rmse --val examples/restaurants/test.csv --indis 500 --gens 50 --target visitors_log1p --seed 5
```

### Boston house prices

```sh
python examples/boston/create_datasets.py
go run main.go fit ./examples/boston/train.csv --val examples/boston/test.csv --loss mae --seed 42 --indis 500 --gens 100
```

### gplearn polynomial example

This example is derived from [gplearn's documentation](http://gplearn.readthedocs.io/en/stable/examples.html). The function to approximate is a trivial and is only here for show.

```sh
python examples/gplearn/create_datasets.py
go run main.go fit examples/gplearn/train.csv --loss mae
koza score examples/gplearn/test.csv --eval mae
```

### Titanic survivors

The data munging is adapted from this [Kaggle kernel](https://www.kaggle.com/scirpus/genetic-programming-lb-0-88).

```sh
python examples/titanic/create_datasets.py
go run .\main.go fit .\examples\titanic\train.csv --loss accuracy --val .\examples\titanic\val.csv --target Survived --ignore PassengerId --parsimony 0.0001 --gens 100 --indis 2000 --funcs sum,sub,mul,div,cos,sin,min,max,pow
go run .\main.go predict examples/titanic/test.csv --output examples/titanic/submission.csv --keep PassengerId --target Survived
```
