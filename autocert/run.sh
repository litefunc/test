#!/bin/sh

docker stop autocert
docker rm autocert
docker rmi autocert:latest
docker load --input autocert.tar

docker run -d --name autocert -p 80:80 -p 443:443 \
    -v /var/lib/.cache/golang-autocert:/var/lib/.cache/golang-autocert \
    --restart always autocert:latest
