# xgp CLI

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

This will install the xgp package together with the CLI. The following Go dependencies will be installed alongside:

- [fatih/color](https://github.com/fatih/color)
- [gosuri/uiprogress](https://github.com/gosuri/uiprogress)
- [MaxHalford/gago](https://github.com/MaxHalford/gago)
- [urfave/cli](https://github.com/urfave/cli)

## Development

```sh
cd xgp/cmd
go run *.go fit ../examples/regression/train.csv -tc y && go run *.go predict ../examples/regression/test.csv -tc y
```
