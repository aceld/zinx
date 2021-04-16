#!/bin/bash

set -e

APP_NAME=$1
APP_VERSION=v$(cat version)
BUILD_VERSION=$(git log -1 --oneline)
BUILD_TIME=$(date "+%FT%T%z")
GIT_REVISION=$(git rev-parse --short HEAD)
GIT_BRANCH=$(git name-rev --name-only HEAD)
GO_VERSION=$(go version)


go build -ldflags " \
	-X 'main.AppName=${APP_NAME}' 			\
	-X 'main.AppVersion=${APP_VERSION}'     \
	-X 'main.BuildVersion=${BUILD_VERSION//\'/_}' \
	-X 'main.BuildTime=${BUILD_TIME}'       \
	-X 'main.GitRevision=${GIT_REVISION}'   \
	-X 'main.GitBranch=${GIT_BRANCH}'       \
	-X 'main.GoVersion=${GO_VERSION}'       \
	" -o $APP_NAME
