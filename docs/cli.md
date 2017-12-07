# Command-line interface (CLI)

## Installation

### Mac

First, [install Go](https://golang.org/dl/), set your `GOPATH`, and make sure `$GOPATH/bin` is on your `PATH`.

```sh
brew install go # If using homebrew

# Put these in your .bash_profile or .zshrc
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

Now you can install the `koza` command-line interface with the following command.

```sh
go install github.com/MaxHalford/koza/cmd/koza
```

Run `koza --help` to check if the installation was successful.

### Linux

### Windows

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


