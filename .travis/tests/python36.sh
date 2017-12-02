#!/bin/bash

cd python-package
conda create -n python36 python=3.6
source activate python36
pip install numpy pytest sklearn
pytest tests/
