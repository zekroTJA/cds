#!/bin/bash

DATE=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
TAG=$(git describe --tags)
COMMIT=$(git rev-parse HEAD)

if [ ! -d builds ]; then
    mkdir builds
fi

echo "Building linux_arm..."
(env GOOS=linux GOARCH=amd64 \
    go build -o builds/cds_linux_arm \
    -ldflags "-X main.appDate=$DATE -X main.appVersion=$TAG -X main.appCommit=$COMMIT")

echo "Building linux_amd64.."
(env GOOS=linux GOARCH=arm \
    go build -o builds/cds_linux_amd64 \
    -ldflags "-X main.appDate=$DATE -X main.appVersion=$TAG -X main.appCommit=$COMMIT")

echo "Building win_amd64..."
(env GOOS=windows GOARCH=amd64 \
    go build -o builds/cds_win_amd64.exe \
    -ldflags "-X main.appDate=$DATE -X main.appVersion=$TAG -X main.appCommit=$COMMIT")

echo "Building mac_amd64..."
(env GOOS=darwin GOARCH=amd64 \
    go build -o builds/cds_mac_amd64 \
    -ldflags "-X main.appDate=$DATE -X main.appVersion=$TAG -X main.appCommit=$COMMIT")

wait