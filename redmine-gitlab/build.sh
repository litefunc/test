#!/bin/bash
BUILD_PATH=$(
	cd $(dirname $0)
	pwd
)
cd $BUILD_PATH
TAG=latest
REPO=redmine-gitlab

export GOARCH=$1
export CGO_ENABLED=0

if [ "$GOARCH" == "" ]; then
	GOARCH=$(go env | grep 'GOARCH=' | sed 's/GOARCH=//g' | sed 's/"//g')
fi

echo $GOARCH

go build -o main.bin test/redmine-gitlab

if [ "$GOARCH" == "amd64" ]; then
	dockerfile="-f $BUILD_PATH/Dockerfile"
	echo $dockerfile
fi

case ${2} in
"img")
	docker build -t $REPO:$TAG $dockerfile $BUILD_PATH
	;;
"tar")
	docker build -t $REPO:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/$REPO.tar $REPO:$TAG
	echo "Image save as : $BUILD_PATH/$REPO.tar"
	docker rmi $REPO:$TAG
	chown $USER:$USER $BUILD_PATH/$REPO.tar
	chmod 777 $BUILD_PATH/$REPO.tar
	;;
"default")
	TAG=default
	docker build -t $REPO:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/$REPO.tar $REPO:$TAG
	echo "Image save as : $BUILD_PATH/$REPO.tar"
	docker rmi $REPO:$TAG
	chown $USER:$USER $BUILD_PATH/$REPO.tar
	chmod 777 $BUILD_PATH/$REPO.tar
	;;
"all")
	docker build -t $REPO:$TAG $dockerfile $BUILD_PATH
	docker save -o $BUILD_PATH/$REPO.tar $REPO:$TAG
	echo "Image save as : $BUILD_PATH/$REPO.tar"
	chown $USER:$USER $BUILD_PATH/$REPO.tar
	chmod 777 $BUILD_PATH/$REPO.tar
	;;
*) ;;

esac
