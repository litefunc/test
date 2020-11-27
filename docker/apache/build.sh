#!/bin/bash
BUILD_PATH=$(
	cd $(dirname $0)
	pwd
)
cd $BUILD_PATH

export GOARCH=$1
export CGO_ENABLED=0

if [ "$GOARCH" == "amd64" ]; then
	dockerfile="-f $BUILD_PATH/Dockerfile_amd64"
	echo $dockerfile
fi

TAG=2.4

case ${2} in
"img")
	docker build -t apache:$TAG $dockerfile $BUILD_PATH
	;;
"tar")
	docker build -t apache:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/apache.tar apache:$TAG
	echo "Image save as : $BUILD_PATH/apache.tar"
	docker rmi apache:$TAG
	chown $USER:$USER $BUILD_PATH/apache.tar
	chmod 777 $BUILD_PATH/apache.tar
	;;
"default")
	TAG=default
	docker build -t apache:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/apache.tar apache:$TAG
	echo "Image save as : $BUILD_PATH/apache.tar"
	docker rmi apache:$TAG
	chown $USER:$USER $BUILD_PATH/apache.tar
	chmod 777 $BUILD_PATH/apache.tar
	;;
"all")
	docker build -t apache:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/apache.tar apache:$TAG
	echo "Image save as : $BUILD_PATH/apache.tar"
	chown $USER:$USER $BUILD_PATH/apache.tar
	chmod 777 $BUILD_PATH/apache.tar
	;;
*) ;;

esac
