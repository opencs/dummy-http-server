#!/bin/bash
go mod vendor
docker run --dns=8.8.8.8 --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp  docker-go-build bash build.sh
