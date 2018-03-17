# Contributing

## To do

See [here](https://github.com/MaxHalford/xgp/projects).

## Documentation

The documentation is rendered with [MkDocs](http://www.mkdocs.org/). Run the following commands to run the documentation locally.

```sh
pip install mkdocs mkdocs-material pygments
cd /path/to/xgp
mkdocs serve
```

## Performance

```sh
go test -bench . --cpuprofile=cpu.prof
go tool pprof --pdf xgp.test cpu.prof > profile.pdf
```

## Magic numbers

Most of the behavior of xgp can be determined by the user. However for various reasons some fields/numbers have to be hard-coded. When this is done the habit I took is to annotate the line of code with a `// MAGIC` comment. In some cases this is due to bad design and should be fixed.

## Steps for adding a new parameter

### Steps

1. Add to the `Estimator`'s fields if needed
2. Add to the `New` method
3. Add to each language package
4. Add to the documentation

### Guidelines

- Order parameters alphabetically
- Respect language-specific conventions for naming
