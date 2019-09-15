#!/bin/sh

# CGO_ENABLED=0 is requried for alpine image
CGO_ENABLED=0 go build -o main.bin