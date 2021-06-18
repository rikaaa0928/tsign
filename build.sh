#!/bin/bash
GOOS=windows GOARCH=amd64 go build -o tsign.exe
GOOS=linux GOARCH=amd64 go build -o tsign
GOOS=darwin GOARCH=amd64 go build -o tsign_darwin
