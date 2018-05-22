#!/bin/bash

go get github.com/mattn/goveralls
go get `go list ./... | grep -v -e cmd`
go test -race -cover -coverprofile=coverage.out `go list ./... | grep -v -e cmd`
${HOME}/gopath/bin/goveralls -coverprofile=coverage.out -service=travis-ci
