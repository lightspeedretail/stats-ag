#!/bin/bash

echo "Building v$1"
go build -o bin/stats-ag -ldflags "-X main.BUILD_DATE `date +%Y-%m-%d` -X main.VERSION $1 -X main.COMMIT_SHA `git rev-parse --verify HEAD`" 
