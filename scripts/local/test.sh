#!/usr/bin/env bash

export GOPATH="$PWD"

go test -test.v -run "$1"  /Users/oleg/codes/huru/src/huru/utils/utils_test.go