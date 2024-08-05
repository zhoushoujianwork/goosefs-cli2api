#!/bin/bash
# build for linux&windows&mac
RELEASE="v0.0.1-beta.$(date +%Y%m%d)"

# GOOS=darwin GOARCH=arm64 go build -o ./bin/goosefs-cli2api-darwin-arm64 -ldflags "-X main.version=$RELEASE"
# GOOS=darwin GOARCH=amd64 go build -o ./bin/goosefs-cli2api-darwin-amd64 -ldflags "-X main.version=$RELEASE"
GOOS=linux GOARCH=amd64 go build -o ./bin/goosefs-cli2api-linux-amd64 -ldflags "-X main.version=$RELEASE"
# GOOS=linux GOARCH=arm64 go build -o ./bin/goosefs-cli2api-linux-arm64 -ldflags "-X main.version=$RELEASE"
# GOOS=windows GOARCH=amd64 go build -o ./bin/goosefs-cli2api-windows-amd64.exe -ldflags "-X main.version=$RELEASE"

echo "build success"