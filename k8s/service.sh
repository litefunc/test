#!/bin/sh

kubectl expose deployment/test --type="NodePort" --port 8088