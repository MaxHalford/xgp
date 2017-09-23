# xgp command-line interface (CLI)

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

The following Go dependencies will be installed alongside:

- [fatih/color](https://github.com/fatih/color)
- [gosuri/uiprogress](https://github.com/gosuri/uiprogress)
- [MaxHalford/gago](https://github.com/MaxHalford/gago)
- [spf13/cobra](https://github.com/spf13/cobra)
- [spf13/viper](https://github.com/spf13/viper)

## Development

```sh
cd xgp/cmd
go run main.go fit ../examples/regression/train.csv && go run main.go predict ../examples/regression/test.csv
```
