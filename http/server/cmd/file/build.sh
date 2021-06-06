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
if [ "$GOARCH" = "arm64" ]; then
    
    export CC=aarch64-linux-gnu-gcc
fi

echo $GOARCH

go build -o main.$GOARCH.bin .

if [ "$GOARCH" == "amd64" ]; then
    dockerfile="-f $BUILD_PATH/Dockerfile_amd64"
    echo $dockerfile
fi

NAME=file_server
TAG=latest

case ${2} in
    "img")
        docker build -t $NAME:$TAG $dockerfile $BUILD_PATH
    ;;
    "tar")
        docker build -t $NAME:$TAG $dockerfile $BUILD_PATH
        docker save -o $BUILD_PATH/$NAME.tar $NAME:$TAG
        echo "Image save as : $BUILD_PATH/$NAME.tar"
        docker rmi $NAME:$TAG
        chown $USER:$USER $BUILD_PATH/$NAME.tar
        chmod 777 $BUILD_PATH/$NAME.tar
    ;;
    "default")
        TAG=default
        docker build -t $NAME:$TAG $dockerfile $BUILD_PATH
        docker save -o $BUILD_PATH/$NAME.tar $NAME:$TAG
        echo "Image save as : $BUILD_PATH/$NAME.tar"
        docker rmi $NAME:$TAG
        chown $USER:$USER $BUILD_PATH/$NAME.tar
        chmod 777 $BUILD_PATH/$NAME.tar
    ;;
    "all")
        docker build -t $NAME:$TAG $dockerfile $BUILD_PATH
        docker save -o $BUILD_PATH/$NAME.tar $NAME:$TAG
        echo "Image save as : $BUILD_PATH/$NAME.tar"
        chown $USER:$USER $BUILD_PATH/$NAME.tar
        chmod 777 $BUILD_PATH/$NAME.tar
    ;;
    *) ;;
    
esac
