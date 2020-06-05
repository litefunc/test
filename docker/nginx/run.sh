#!/bin/sh

FILE_DIR=$(
    cd $(dirname $0)
    pwd
)

# docker run --name web -p 8080:80 \
#     -v $FILE_DIR/nginx:/etc/nginx nginx:test

# docker run --rm --name web -p 8080:80 \
#     -v $FILE_DIR/ssl:/etc/nginx/ssl nginx:test

# docker run -d --name web -p 80:80 -p 443:443 -p 8080:8080 \
#     -v $FILE_DIR/ssl:/etc/nginx/ssl nginx:test

# docker run -d --name web --network host \
#     -v $FILE_DIR/ssl:/etc/nginx/ssl nginx:test

docker run -d --name web --network host \
    -v $FILE_DIR/ssl:/etc/nginx/ssl nginx:test

docker run -d --name web --network host nginx:test
