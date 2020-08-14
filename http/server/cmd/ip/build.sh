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

go build -o main.$GOARCH.bin -ldflags="-X 'main.version=$GIT_DATE'" test/http/cmd/ip

dockerfile="-f $BUILD_PATH/Dockerfile"

case ${2} in
"img")
	docker build -t ip:$TAG $dockerfile $BUILD_PATH
	;;
"tar")
	docker build -t ip:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/ip.tar ip:$TAG
	echo "Image save as : $BUILD_PATH/ip.tar"
	docker rmi ip:$TAG
	chown $USER:$USER $BUILD_PATH/ip.tar
	chmod 777 $BUILD_PATH/ip.tar
	;;
"default")
	TAG=default
	docker build -t ip:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/ip.tar ip:$TAG
	echo "Image save as : $BUILD_PATH/ip.tar"
	docker rmi ip:$TAG
	chown $USER:$USER $BUILD_PATH/ip.tar
	chmod 777 $BUILD_PATH/ip.tar
	;;
"all")
	docker build -t ip:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/ip.tar ip:$TAG
	echo "Image save as : $BUILD_PATH/ip.tar"
	chown $USER:$USER $BUILD_PATH/ip.tar
	chmod 777 $BUILD_PATH/ip.tar
	;;
*) ;;

esac
