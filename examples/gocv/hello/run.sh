#!/usr/bin/env bash

# STEP 1: Determinate the required values

PACKAGE="github.com/pubgo/funk"
VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

# STEP 2: Build the ldflags

LDFLAGS=(
  "-s -w"
  "-X '${PACKAGE}/version.Version=${VERSION}'"
  "-X '${PACKAGE}/version.CommitHash=${COMMIT_HASH}'"
  "-X '${PACKAGE}/version.BuildTime=${BUILD_TIMESTAMP}'"
)

# STEP 3: Actual Go build process

# brew install musl-cross
GO_CC=aarch64-linux-musl-gcc
localos=$(uname)
if [ "$localos" == "Linux" ]; then
  #in linux use https://armkeil.blob.core.windows.net/developer/Files/downloads/gnu/13.3.rel1/binrel/arm-gnu-toolchain-13.3.rel1-x86_64-arm-none-linux-gnueabihf.tar.xz
  GO_CC=arm-none-linux-gnueabihf-gcc
fi

CGO_ENABLED=1 GOARCH=arm64 GOOS=linux CC=${GO_CC} CGO_LDFLAGS="-static" go build -v -gcflags '-N -l' -ldflags="${LDFLAGS[*]}" -a  -installsuffix cgo -o manager
