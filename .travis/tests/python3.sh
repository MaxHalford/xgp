#!/bin/bash

if [ ${TRAVIS_OS_NAME} == "osx" ]; then
    wget -O conda.sh https://repo.continuum.io/miniconda/Miniconda3-latest-MacOSX-x86_64.sh
else
    wget -O conda.sh https://repo.continuum.io/miniconda/Miniconda3-latest-Linux-x86_64.sh
fi
bash conda.sh -b -p $HOME/miniconda
export PATH="$HOME/miniconda/bin:$PATH"
hash -r
conda config --set always_yes yes --set changeps1 no
conda update -q conda
conda create -n python3 python=3.5

cd python-package
source activate python3
python -m pip install numpy pytest sklearn
py.test tests/
