#!/bin/bash

# Install the Python dependencies needed to generate the datasets
if [ ${TRAVIS_OS_NAME} == "osx" ]; then
    wget -O miniconda.sh https://repo.continuum.io/miniconda/Miniconda3-latest-MacOSX-x86_64.sh
else
    wget -O miniconda.sh https://repo.continuum.io/miniconda/Miniconda3-latest-Linux-x86_64.sh
fi
bash miniconda.sh -b -p $HOME/miniconda
export PATH="$HOME/miniconda/bin:$PATH"
hash -r
conda config --set always_yes yes --set changeps1 no
conda update -q conda
pip install scipy pandas scikit-learn

# Naviguate to the CLI root directory
cd cmd/xgp

# Boston
python examples/boston/create_datasets.py
go run main.go fit examples/boston/train.csv \
    --loss mae \
    --indis 20 \
    --gens 10 \
    --output examples/boston/ensemble.json \
    --seed 42
go run main.go predict examples/boston/test.csv \
    --output examples/boston/predictions.csv \
    --ensemble examples/boston/ensemble.json

# Breast cancer
python examples/breast_cancer/create_datasets.py
go run main.go fit examples/breast_cancer/train.csv \
    --loss logloss \
    --eval accuracy \
    --val examples/breast_cancer/val.csv \
    --target Survived \
    --ignore PassengerId \
    --parsimony 0.001 \
    --gens 30 \
    --indis 100 \
    --funcs sum,sub,mul,div,cos,sin,min,max
go run main.go predict examples/breast_cancer/val.csv \
    --output examples/breast_cancer/submission.csv \
    --target has_cancer \
    --proba
