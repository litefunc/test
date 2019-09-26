#!/bin/sh
BUILD_PATH=$(cd $(dirname $0); pwd)
cd $BUILD_PATH

GIT_HASH=$(git log -1 --format=%h)
TAG=latest


build(){
	docker build -t test:$TAG $BUILD_PATH
}



case $1 in
	"") build;;
	*)   build;;
esac