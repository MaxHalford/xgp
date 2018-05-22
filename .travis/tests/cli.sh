#!/bin/bash

# Install the Python dependencies needed to generate the datasets
pip install scipy pandas scikit-learn

# Naviguate to the CLI root directory
cd cmd/xgp

# Boston bagging
python examples/boston/create_datasets.py
go run main.go fit examples/boston/train.csv \
    --mode bagging \
    --loss mae \
    --indis 20 \
    --gens 10 \
    --output examples/boston/ensemble.json \
    --seed 42
go run main.go predict examples/boston/test.csv \
    --output examples/boston/predictions.csv \
    --ensemble examples/boston/ensemble.json

# Breast cancer bagging
python examples/breast_cancer/create_datasets.py
go run main.go fit examples/breast_cancer/train.csv \
    --mode bagging \
    --loss logloss \
    --eval accuracy \
    --val examples/breast_cancer/val.csv \
    --target has_cancer \
    --ignore PassengerId \
    --parsimony 0.001 \
    --gens 30 \
    --indis 100 \
    --funcs sum,sub,mul,div,cos,sin,min,max
go run main.go predict examples/breast_cancer/val.csv \
    --output examples/breast_cancer/submission.csv \
    --target has_cancer \
    --proba
