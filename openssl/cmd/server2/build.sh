#!/bin/bash
BUILD_PATH=$(
	cd $(dirname $0)
	pwd
)
cd $BUILD_PATH
GIT_DATE=$(git log -1 --format=%cd)
TAG=$(git log -1 --format=%h)

export GOARCH=$1
export CGO_ENABLED=0

if [ "$GOARCH" == "" ]; then
	GOARCH=$(go env | grep 'GOARCH=' | sed 's/GOARCH=//g' | sed 's/"//g')
fi

echo $GOARCH

go build -o main.$GOARCH.bin test/openssl/cmd/server2

if [ "$GOARCH" == "amd64" ]; then
	dockerfile="-f $BUILD_PATH/Dockerfile_amd64"
	echo $dockerfile
fi

case ${2} in
"img")
	docker build -t openssl:$TAG $dockerfile $BUILD_PATH
	;;
"tar")
	docker build -t openssl:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/openssl.tar openssl:$TAG
	echo "Image save as : $BUILD_PATH/openssl.tar"
	docker rmi openssl:$TAG
	chown $USER:$USER $BUILD_PATH/openssl.tar
	chmod 777 $BUILD_PATH/openssl.tar
	;;
"default")
	TAG=default
	docker build -t openssl:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/openssl.tar openssl:$TAG
	echo "Image save as : $BUILD_PATH/openssl.tar"
	docker rmi openssl:$TAG
	chown $USER:$USER $BUILD_PATH/openssl.tar
	chmod 777 $BUILD_PATH/openssl.tar
	;;
"all")
	docker build -t openssl:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/openssl.tar openssl:$TAG
	echo "Image save as : $BUILD_PATH/openssl.tar"
	chown $USER:$USER $BUILD_PATH/openssl.tar
	chmod 777 $BUILD_PATH/openssl.tar
	;;
*) ;;

esac
