#!/bin/sh


export GOARCH=$1
if [ "$GOARCH" = "arm64" ]; then
	export CGO_ENABLED=1
	export CC=aarch64-linux-gnu-gcc
fi

go build -o httpserver .

