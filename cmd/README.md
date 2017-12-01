# Command-line interface (CLI)

## Development

```sh
cd koza/cmd
go run main.go fit ../examples/gplearn/train.csv --loss_metric mae
go run main.go predict ../examples/gplearn/test.csv -e mae
```

```sh
cd koza/cmd
go run main.go fit ../examples/boston/train.csv --loss mae # 6.754572271386434 with median
go run main.go predict ../examples/boston/test.csv --eval mae # 6.119760479041916 with median
```

```sh
cd koza/cmd
go run main.go fit ../examples/iris/train.csv -l accuracy
go run main.go predict ../examples/iris/test.csv -e accuracy
```


```sh
cd koza/cmd
go run main.go fit ../examples/titanic/train.csv --loss f1_score --target Survived
go run main.go predict ../examples/titanic/val.csv --eval f1_score --target Survived
```
