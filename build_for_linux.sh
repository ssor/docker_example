#!/bin/bash
cp config.yaml ~/docker_vol/
CGO_ENABLED=0 GOOS=linux  go build -o hello main.go