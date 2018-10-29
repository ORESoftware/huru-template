#!/usr/bin/env bash

export GOPATH="$PWD"
export GOCACHE="on"  # on / off

if [[ -z "$2" ]]; then
     go test -test.v "$1" 
else
     go test -test.v -run "$1"  "$2"
fi

