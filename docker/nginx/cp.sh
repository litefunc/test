#!/bin/sh

FILE_DIR=$(
    cd $(dirname $0)
    pwd
)

# docker cp $FILE_DIR/default.conf web:/etc/nginx/conf.d/default.conf

# docker cp $FILE_DIR/nginx.conf web:/etc/nginx/nginx.conf

# docker cp $FILE_DIR/nginx.conf web:/etc/nginx/conf.d/nginx.conf

# docker cp $FILE_DIR/nginx.conf web:/etc/nginx/conf.d/nginx.conf

# docker cp $FILE_DIR/html/index.html web:/usr/share/nginx/html/index.html

# docker cp $FILE_DIR/html/index1.html web:/usr/share/nginx/html/web/index.html

# docker cp $FILE_DIR/web web:/usr/share/nginx/html/web

# docker cp index.html demo:/etc/nginx/conf.d/nginx.conf

# docker cp html/index.html nginx:/root/nginx/html/

docker cp html/index1.html nginx:/usr/share/nginx/html/

docker cp $FILE_DIR/nginx.conf nginx:/etc/nginx/conf.d/nginx.conf
