#!/bin/bash

if [ ${TASK} == "go_test" ]; then
    set -e
    go get github.com/mattn/goveralls
    go test `go list ./... | grep -v -e cmd -e python-package`
    ${HOME}/gopath/bin/goveralls -coverprofile=coverage.out -service=travis-ci
fi
