#!/bin/bash
# build for linux&windows&mac
RELEASE="v0.0.1-beta.$(date +%Y%m%d)"

# GOOS=darwin GOARCH=arm64 go build -o ./bin/goosefs-cli2api-darwin-arm64 -ldflags "-X main.version=$RELEASE"
# GOOS=darwin GOARCH=amd64 go build -o ./bin/goosefs-cli2api-darwin-amd64 -ldflags "-X main.version=$RELEASE"
GOOS=linux GOARCH=amd64 go build -o ./bin/goosefs-cli2api-linux-amd64 -ldflags "-X main.version=$RELEASE"
# GOOS=linux GOARCH=arm64 go build -o ./bin/goosefs-cli2api-linux-arm64 -ldflags "-X main.version=$RELEASE"
# GOOS=windows GOARCH=amd64 go build -o ./bin/goosefs-cli2api-windows-amd64.exe -ldflags "-X main.version=$RELEASE"

echo "build success"


push() {
    docker build -t zhoushoujian/goosefs-cli2api:$RELEASE . --build-arg="RELEASE=$RELEASE"
    if [ $? -eq 0 ]; then
        docker push zhoushoujian/goosefs-cli2api:$RELEASE
        docker tag zhoushoujian/goosefs-cli2api:$RELEASE zhoushoujian/goosefs-cli2api:latest
        docker push zhoushoujian/goosefs-cli2api:latest
    else
        echo "build failed"
    fi
}

if [ "$1" = "push" ]; then
    echo "push"
    push
fi