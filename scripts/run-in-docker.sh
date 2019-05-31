#!/bin/bash

docker run -e GIT_SSL_NO_VERIFY=1 --rm -v "$PWD":/usr/src/dhara -v /Users/mangarg/docker-go-mod:/go -w /usr/src/dhara -p 5678:5678 -it golang:1.12 bash scripts/buildrun.sh