#!/bin/bash -x
mkdir -p $GOPATH/{src,pkg,bin}
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/craigdfrench/event
rm /tmp/*.pid
$GOPATH/src/github.com/craigdfrench/event/cmd/event-service/event-service.sh build
$GOPATH/src/github.com/craigdfrench/event/cmd/event-service/event-service.sh start
