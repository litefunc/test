#!/bin/sh
BUILD_PATH=$(
	cd $(dirname $0)
	pwd
)
cd $BUILD_PATH

GIT_HASH=$(git log -1 --format=%h)
TAG=test

build() {
	docker build -t nginx:$TAG $BUILD_PATH
}

case $1 in
"") build ;;
*) build ;;
esac
