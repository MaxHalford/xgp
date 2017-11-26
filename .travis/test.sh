#!/bin/bash

if [ ${TASK} == "test_core" ]; then
    set -e
    go get github.com/mattn/goveralls
    go test -race -cover -coverprofile=coverage.out `go list ./... | grep -v -e cmd -e python-package`
    ${HOME}/gopath/bin/goveralls -coverprofile=coverage.out -service=travis-ci
fi

if [ ${TASK} == "test_python" ]; then
    set -e
    cd python-package
    pytest
fi
