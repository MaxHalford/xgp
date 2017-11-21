# Command-line interface (CLI)

## Development

```sh
cd xgp/cmd
go run main.go fit ../examples/gplearn/train.csv --loss_metric mae
go run main.go predict ../examples/gplearn/test.csv -e mae
```

```sh
cd xgp/cmd
go run main.go fit ../examples/iris/train.csv -l accuracy
go run main.go predict ../examples/iris/test.csv -e accuracy
```


```sh
cd xgp/cmd
go run main.go fit ../examples/titanic/train.csv -l f1_score -y Survived
go run main.go predict ../examples/titanic/val.csv -e f1_score
```
