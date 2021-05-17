#!/bin/bash
BUILD_PATH=$(
	cd $(dirname $0)
	pwd
)
cd $BUILD_PATH


export GOARCH=$1
export CGO_ENABLED=0

if [ "$GOARCH" == "" ]; then
	GOARCH=$(go env | grep 'GOARCH=' | sed 's/GOARCH=//g' | sed 's/"//g')
fi

echo $GOARCH

go build -o udp-receiver .

