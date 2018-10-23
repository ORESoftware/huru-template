#!/usr/bin/env bash
export GOPATH="$(cd $(dirname "$BASH_SOURCE") && pwd)"
go get "$1"