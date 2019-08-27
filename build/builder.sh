#!/bin/bash

if [ "$1" = 't' ] || [ "$1" = 'test' ]; then
     GOOS=linux GOARCH=arm GOARM=5 go build ./TestBranch/test.go 
fi

if [ "$1" = 'm' ] || [ "$1" = 'main' ]; then
     GOOS=linux GOARCH=arm GOARM=5 go build ./cmd/main.go 
fi