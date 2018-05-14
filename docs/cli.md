# Command-line interface (CLI)

## Installation

First, [install Go](https://golang.org/dl/), set your `GOPATH`, and make sure `$GOPATH/bin` is on your `PATH`. Here are some additional ressources depending on your operating system:

- [Linux](https://www.tecmint.com/install-go-in-linux/)
- [Mac](http://sourabhbajaj.com/mac-setup/Go/README.html)
- [Windows](http://www.wadewegner.com/2014/12/easy-go-programming-setup-for-windows/)

Next, regardless of your OS, you can install the `xgp` CLI with the following command.

```sh
go install github.com/MaxHalford/xgp/cmd/xgp
```

If `xgp --help` runs without any errors then the installation was successful. If you encounter an error feel free to [open an issue on GitHub](https://github.com/MaxHalford/xgp/issues/new).

## Usage

!!! tip
    Apart from the following documentation you can also check out [the command-line usage examples](https://github.com/MaxHalford/xgp/tree/master/cmd/xgp#examples).

!!! tip
    Run `xgp <command> -h` to get help with a command. For example `xgp fit -h` will display the help for the `fit` command.

### Training

The `fit` command trains programs against a training dataset and saves the best one to a JSON file. The only required argument is a path to a CSV file which acts as the training dataset. The dataset is that it should contain only numerical data. Moreover the first row should contain column names.

!!! warning
    For the while XGP does not handle categorical data. You should preemptively encode the categorical features in your dataset before feeding it to XGP. The recommended way is to use [label encoding](http://scikit-learn.org/stable/modules/preprocessing_targets.html#label-encoding) for ordinal data and [one-hot encoding](http://scikit-learn.org/stable/modules/preprocessing.html#encoding-categorical-features) for non-ordinal data.

!!! warning
    For the while xgp does not handle missing values.

Once your dataset is ready, you can train XGP on it with the following command.

```sh
>>> xgp fit train.csv
```

This will evaluate and evolve many programs with default values before finally outputting the best obtained program to a JSON file. By default the JSON file is named `program.json`. The JSON file can then be used by the `predict` command to make predictions on another dataset.

!!! info
    Whether the task is classification or regression is guessed from the loss metric parameter. The available loss metrics are listed [here](training-parameters.md#loss-metrics)

There are many parameters you can use; the details and default values and are specified in the [training parameters section](training-parameters.md)

```sh
>>> xgp fit train.csv --loss mse --val mae --gens 64 --indis 256 --parsimony 0.001
```

In addition to these parameters, the following ones are available with the `fit` command.

| Argument | Description | Default |
|----------|-------------|---------|
| ignore | Comma-separated list of columns to ignore | |
| output | Path where to save the JSON representation of the best program | `program.json` |
| target | Name of the target column in the training and validation datasets | `y` |
| val | Path to a validation dataset that can be used to monitor out-of-bag performance | |
| verbose | Indicates whether or not progress should be displayed | `true` |

If you use the `val` argument then the best model of each generation will be scored against the validation dataset. The resulting score is called the out-of-bag score because it is obtained by making predictions on a dataset that the model hasn't seen.

### Predicting

Once you have produced a program with the `fit` command you can use it to make predictions on a dataset. The test set should have exactly the same format as the training set. Specifically **the columns in the test set should be ordered in the same way they were in the training set**.

```sh
>>> xgp predict test.csv
```

This will make predictions on `test.csv` and save them to a specificied path. The default path is `y_pred.csv`. The following arguments are available for the `predict` command.

| Argument | Description | Default |
|----------|-------------|---------|
| keep | Comma-separated list of columns to keep in the CSV output | |
| output | Path to the CSV output | y_pred.csv |
| program | Path to the program used to make predictions | program.json |
| target | Name of the target column in the CSV output | y |


### Scoring

If you don't want to save predictions and instead only want to evaluate a program then you can use the `score` command. The `score` command will open a program, make predictions against a given dataset, and output a prediction score. By default the scoring metric is the loss metric used for training.

```sh
>>> xgp score test.csv
```

The following arguments are available for the `score` command.

| Argument | Description | Default |
|----------|-------------|---------|
| eval | Evaluation metric | Same as the loss metric used during training |
| program | Path to the program to score | program.json |
| target | Name of the target column in the dataset | y |


### Visualization

Because programs can be represented as trees, [Graphviz](https://www.graphviz.org/) can be used to visualize them. The `todot` command takes a program as input and outputs the Graphviz representation of the program. You can then copy/paste the output and use a service such as [webgraphviz](http://www.webgraphviz.com/) to obtain the visualization. By default the output will not be saved to a file but will however be displayed in the terminal.

```sh
>>> xgp todot program.json
```

You can also feed the `todot` command a formula instead of a JSON file.

```sh
>>> xgp todot "sum(X[13], 42)"
```

The following arguments are available for the `todot` command.

| Argument | Description | Default |
|----------|-------------|---------|
| output | Path to the DOT file output | program.dot |
| save | Save to a DOT file or not | False |
| shell | Output in the terminal or not | True |
