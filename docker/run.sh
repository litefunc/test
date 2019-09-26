#!/bin/sh

FILE_DIR=$(cd $(dirname $0); pwd)

# docker run --rm -it test:latest sh

# docker run --rm -it -v $FILE_DIR/config.json:/root/config/config.json  test:latest sh

# docker run --rm -it -v $FILE_DIR/config:/root/config  test:latest sh

docker run --rm -it -v $FILE_DIR/var/lib/.cache/golang-autocert:/var/lib/.cache/golang-autocert  test:latest sh

