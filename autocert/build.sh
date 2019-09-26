#!/bin/sh
BUILD_PATH=$(
	cd $(dirname $0)
	pwd
)
cd $BUILD_PATH

GIT_HASH=$(git log -1 --format=%h)
TAG=latest
NAME=autocert

export CGO_ENABLED=0
go build .

build() {
	docker build -t $NAME:$TAG $BUILD_PATH
}

tar() {
	build
	docker save -o $BUILD_PATH/$NAME.tar $NAME:$TAG
	echo "Image save as : $BUILD_PATH/$NAME.tar"
	docker rmi $NAME:$TAG
	chown $USER:$USER $BUILD_PATH/$NAME.tar
	chmod 777 $BUILD_PATH/$NAME.tar
}

case $1 in
"") build ;;
"tar") tar ;;
*) build ;;
esac
