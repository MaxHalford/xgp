#!/bin/bash

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
pip install mkdocs==1.0.1 mkdocs-material==3.0.3 python-markdown-math==0.6
cd ..
mkdocs build --verbose --clean --strict
