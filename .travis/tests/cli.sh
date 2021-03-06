#!/bin/bash

# Install Miniconda
mkdir -p download
cd download
echo "Cached in $HOME/download :"
ls -l
echo
if [[ ! -f miniconda.sh ]]
    then
        if [ ${TRAVIS_OS_NAME} == "osx" ]; then
            wget -O miniconda.sh https://repo.continuum.io/miniconda/Miniconda3-latest-MacOSX-x86_64.sh
        else
            wget -O miniconda.sh https://repo.continuum.io/miniconda/Miniconda3-latest-Linux-x86_64.sh
        fi
    fi
chmod +x miniconda.sh && ./miniconda.sh -b
export PATH=/home/travis/miniconda3/bin:$PATH
pip install scipy pandas scikit-learn
cd ..

# Naviguate to the CLI root directory
cd cmd/xgp

# Boston
python examples/boston/create_datasets.py
go run main.go fit examples/boston/train.csv \
    --loss mse \
    --indis 20 \
    --gens 10 \
    --output examples/boston/model.json \
    --seed 42
go run main.go predict examples/boston/test.csv \
    --output examples/boston/predictions.csv \
    --model examples/boston/model.json

# Breast cancer bagging
python examples/breast_cancer/create_datasets.py
go run main.go fit examples/breast_cancer/train.csv \
    --loss logloss \
    --eval accuracy \
    --val examples/breast_cancer/val.csv \
    --target has_cancer \
    --ignore PassengerId \
    --parsimony 0.001 \
    --gens 30 \
    --indis 100 \
    --funcs add,sub,mul,div,cos,sin,min,max
go run main.go predict examples/breast_cancer/val.csv \
    --output examples/breast_cancer/submission.csv \
    --target has_cancer \
    --proba
