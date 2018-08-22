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
pip install mkdocs mkdocs-material python-markdown-math
cd ..
mkdocs build --verbose --clean --strict
