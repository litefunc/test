#!/bin/bash

docker run -d --name file_server -p 9000:9000 -v `pwd`/static:/static --restart always file_server:latest