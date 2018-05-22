#!/bin/bash

cd
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
export PATH=$HOME/travis/miniconda3/bin:$PATH
pip install scipy pandas scikit-learn
cd ../..

# Naviguate to the CLI root directory
ls
cd xgp/cmd/xgp

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
