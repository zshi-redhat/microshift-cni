#!/usr/bin/env bash

set -eu

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

export BIN_PATH=bin

mkdir -p ${BIN_PATH}

echo "Building microshift ovn-kubernetes cni ..."
CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "-s -w" -o ${BIN_PATH}/ovn-kubernetes-cni main.go
