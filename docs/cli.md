# Command-line interface (CLI)

## Installation

First, [install Go](https://golang.org/dl/), set your `GOPATH`, and make sure `$GOPATH/bin` is on your `PATH`.

```sh
brew install go # If using homebrew

# Put these in your .bash_profile or .zshrc
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

Now you can install the xgp command-line interface with the following command.

```sh
go get -u github.com/MaxHalford/xgp/cmd
```

## Commands


!!! tip
    Whichever the command, you may run `xgp <command> -h` to get help with it.

### fit

#### Basic example

```sh
xgp fit train.csv
```

#### Full example


### predict

#### Basic example

```sh
xgp predict test.csv
```

#### Full example


