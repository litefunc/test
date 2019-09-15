#!/bin/sh

kubectl run test --image=litefunc/simple-httpserver:v1 --port 8088