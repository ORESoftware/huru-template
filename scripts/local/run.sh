#!/usr/bin/env bash

# export GOPATH="$(cd $(dirname "$BASH_SOURCE") && pwd)"
export GOPATH="$PWD";
go install huru
"$GOPATH/bin/huru"