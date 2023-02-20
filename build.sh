#!/usr/bin/env sh
set CGO_ENABLE=1
set GOOS=linux
set GOARCH=amd64

#@rem go build -ldflags "-s -w"
go build -ldflags "-w"