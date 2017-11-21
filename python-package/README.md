# koza Python package

## Installation

```sh
git clone github.com/MaxHalford/koza
cd koza/python-package
go get
go build -buildmode=c-shared -o koza.so *.go
```

## Development

```sh
pip install -e .[dev]
```
