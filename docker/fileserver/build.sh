#!/bin/bash
BUILD_PATH=$(
	cd $(dirname $0)
	pwd
)
cd $BUILD_PATH

TAG=0.0.1

export GOARCH=$1
export CGO_ENABLED=0

if [ "$GOARCH" == "" ]; then
	GOARCH=$(go env | grep 'GOARCH=' | sed 's/GOARCH=//g' | sed 's/"//g')
fi

echo $GOARCH

go build -o main.$GOARCH.bin .

if [ "$GOARCH" == "amd64" ]; then
	dockerfile="-f $BUILD_PATH/Dockerfile_amd64"
	echo $dockerfile
fi

case ${2} in
"img")
	docker build -t file_server:$TAG $dockerfile $BUILD_PATH
	;;
"tar")
	docker build -t file_server:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/file_server.tar file_server:$TAG
	echo "Image save as : $BUILD_PATH/file_server.tar"
	docker rmi file_server:$TAG
	chown $USER:$USER $BUILD_PATH/file_server.tar
	chmod 777 $BUILD_PATH/file_server.tar
	;;
"default")
	TAG=default
	docker build -t file_server:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/file_server.tar file_server:$TAG
	echo "Image save as : $BUILD_PATH/file_server.tar"
	docker rmi file_server:$TAG
	chown $USER:$USER $BUILD_PATH/file_server.tar
	chmod 777 $BUILD_PATH/file_server.tar
	;;
"all")
	docker build -t file_server:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/file_server.tar file_server:$TAG
	echo "Image save as : $BUILD_PATH/file_server.tar"
	chown $USER:$USER $BUILD_PATH/file_server.tar
	chmod 777 $BUILD_PATH/file_server.tar
	;;
*) ;;

esac
