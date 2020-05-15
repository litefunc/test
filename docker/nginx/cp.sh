#!/bin/sh

FILE_DIR=$(
    cd $(dirname $0)
    pwd
)

# docker cp $FILE_DIR/nginx.conf web:/etc/nginx/nginx.conf

docker cp $FILE_DIR/nginx.conf web:/etc/nginx/conf.d/nginx.conf

# docker cp $FILE_DIR/nginx1.conf web:/etc/nginx/conf.d/nginx.conf

# docker cp $FILE_DIR/html web:/etc/nginx/html

docker cp $FILE_DIR/web web:/usr/share/nginx/html
