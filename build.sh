#!/bin/sh
#GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main
#CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags musl -ldflags="-s -w" -o main ./main.go
docker build -t devmgr .