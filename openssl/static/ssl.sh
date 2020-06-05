#!/bin/sh

openssl req -x509 -new -nodes -sha256 -utf8 -days 36500 -newkey rsa:2048 -keyout server.key -out server.crt -config ssl.conf
