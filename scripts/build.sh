#!/bin/bash

TAG=$(git describe --tags)

if [ ! -d bin ]; then
    mkdir ./bin
fi

dep ensure -v

echo "Building linux_arm..."
(env GOOS=linux GOARCH=arm \
    go build -o ./bin/cds_linux_arm \
    -ldflags "-X github.com/zekroTJA/cds/internal/static.AppVersion=$TAG -X github.com/zekroTJA/cds/static.Release=TRUE" \
        ./cmd/cds/main.go)

echo "Building linux_amd64.."
(env GOOS=linux GOARCH=amd64 \
    go build -o ./bin/cds_linux_amd64 \
    -ldflags "-X github.com/zekroTJA/cds/internal/static.AppVersion=$TAG -X github.com/zekroTJA/cds/static.Release=TRUE" \
        ./cmd/cds/main.go)

echo "Building win_amd64..."
(env GOOS=windows GOARCH=amd64 \
    go build -o ./bin/cds_win_amd64.exe \
    -ldflags "-X github.com/zekroTJA/cds/internal/static.AppVersion=$TAG -X github.com/zekroTJA/cds/static.Release=TRUE" \
        ./cmd/cds/main.go)

echo "Building mac_amd64..."
(env GOOS=darwin GOARCH=amd64 \
    go build -o ./bin/cds_mac_amd64 \
    -ldflags "-X github.com/zekroTJA/cds/internal/static.AppVersion=$TAG -X github.com/zekroTJA/cds/static.Release=TRUE" \
        ./cmd/cds/main.go)