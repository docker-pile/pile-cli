#!/bin/bash

mkdir -p build
GOOS=darwin GOARCH=amd64 go build -C src -o ../build/pile-darwin-amd64
shasum -a 256  build/cli-darwin-amd64 

GOOS=darwin GOARCH=arm64 go build -C src -o ../build/pile-darwin-arm64
shasum -a 256  build/cli-darwin-arm64 

GOOS=linux GOARCH=amd64 go build -C src -o ../build/pile-linux-amd64
shasum -a 256  build/cli-linux-amd64
