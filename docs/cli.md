# Command-line interface (CLI)

## Installation

First, [install Go](https://golang.org/dl/), set your `GOPATH`, and make sure `$GOPATH/bin` is on your `PATH`.

```sh
brew install go # If using homebrew

# Put these in your .bash_profile or .zshrc
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

Now you can install the `koza` command-line interface with the following command.

```sh
go get -u github.com/MaxHalford/koza/cmd
```

## Commands


!!! tip
    Run `koza <command> -h` to get help with a command.

### fit

#### Basic example

```sh
koza fit train.csv
```

#### Full example


### predict

#### Basic example

```sh
koza predict test.csv
```

#### Full example


