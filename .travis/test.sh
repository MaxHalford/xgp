#!/bin/bash

if [ ${TASK} == "test_core" ]; then
    set -e
    go get github.com/mattn/goveralls
    go test -race -cover -coverprofile=coverage.out `go list ./... | grep -v -e cmd -e python-package`
    ${HOME}/gopath/bin/goveralls -coverprofile=coverage.out -service=travis-ci
fi

if [ ${TASK} == "test_python_3" ]; then
    set -e
    cd python-package
    source activate python3
    python -m pip install pytest
    py.test tests/
fi
