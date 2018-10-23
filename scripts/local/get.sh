#!/usr/bin/env bash

# export GOPATH="$(cd $(dirname "$BASH_SOURCE") && pwd)"
export GOPATH="$PWD";
go get "$1"